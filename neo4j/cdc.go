package neo4j

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mikeyuniverse/neo4j-cdc-go/entities"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	cdcEarliestQuery = `CALL db.cdc.earliest()`

	logQuery = `
		CALL db.cdc.query($id, [
    		{select: "e"}
    	])
		YIELD id, txId, seq, metadata, event
		RETURN id, txId, seq, metadata, event
		LIMIT $limit`

	showDatabasesWithoutCDC = `
		SHOW DATABASES
		YIELD name, options
		WHERE
		  (options["txLogEnrichment"] IS NULL
		  OR options["txLogEnrichment"] = 'OFF')
		  AND name <> 'system'
		RETURN name`

	changeCDCQuery = `ALTER DATABASE $dbname SET OPTION txLogEnrichment $cdc`
)

func (n *Neo4j) Earliest(ctx context.Context, database string) (string, error) {
	if database == "" {
		return "", errors.New("empty database name")
	}

	s := n.d.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: database,
		FetchSize:    1,
	})

	results, err := s.Run(ctx, cdcEarliestQuery, nil, neo4j.WithTxTimeout(txTimeout))
	if err != nil {
		return "", fmt.Errorf("run query: %w", err)
	}

	earliestID, err := neo4j.SingleTWithContext(ctx, results, func(r *neo4j.Record) (string, error) {
		return stringFromRecord(r, "id")
	})
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	return earliestID, nil
}

func (n *Neo4j) GetCDCItem(ctx context.Context, database string, id string, limit uint) ([]*entities.TxLog, error) {
	if database == "" {
		return nil, errors.New("empty database name")
	}
	if id == "" {
		return nil, errors.New("empty id")
	}
	if limit == 0 {
		return nil, errors.New("limit is 0")
	}

	var logs []*entities.TxLog
	params := map[string]any{
		"id":    id,
		"limit": limit,
	}

	tries := 1

	for {
		s := n.d.NewSession(ctx, neo4j.SessionConfig{
			AccessMode:   neo4j.AccessModeRead,
			DatabaseName: database,
		})

		results, err := s.Run(ctx, logQuery, params, neo4j.WithTxTimeout(txTimeout))
		if err != nil {
			if neo4j.IsRetryable(err) {
				tries++

				slog.Debug(
					"retrying on connectivity error",
					"tries", tries)

				err = s.Close(ctx)
				if err != nil {
					return nil, fmt.Errorf("close session: %w", err)
				}
				continue
			}
			return nil, fmt.Errorf("do query: %w", err)
		}

		logs, err = neo4j.CollectTWithContext(ctx, results, func(r *neo4j.Record) (*entities.TxLog, error) {
			return txLog(r)
		})
		if err != nil {
			return nil, fmt.Errorf("collect logs: %w", err)
		}

		err = s.Close(ctx)
		if err != nil {
			return nil, fmt.Errorf("close session: %w", err)
		}
		break

	}

	return logs, nil
}

// FIXME: update all databases in single query
func (n *Neo4j) EnableFullCDC(ctx context.Context, databases []string) error {
	s := n.d.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	result, err := s.Run(ctx, showDatabasesWithoutCDC, nil, neo4j.WithTxTimeout(txTimeout))
	if err != nil {
		return fmt.Errorf("run query: %w", err)
	}

	databases, err = neo4j.CollectTWithContext(ctx, result, func(r *neo4j.Record) (string, error) {
		return stringFromRecord(r, "name")
	})

	for _, db := range databases {
		result, err := s.Run(ctx, changeCDCQuery, map[string]any{
			"dbname": db,
			"cdc":    "FULL",
		}, neo4j.WithTxTimeout(txTimeout))
		if err != nil {
			return fmt.Errorf("do query for %q: %w", db, err)
		}

		_, err = result.Consume(ctx)
		if err != nil {
			return fmt.Errorf("consume: %w", err)
		}
	}

	return nil
}

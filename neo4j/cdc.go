package neo4j

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	cdcEarliestQuery = `CALL db.cdc.earliest()`
	logQuery         = `CALL db.cdc.query($id)`
)

func (n *Neo4j) Earliest(ctx context.Context) (string, error) {
	s := n.d.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
		FetchSize:  1,
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

type TxLog struct {
}

func (n *Neo4j) LogQuery(ctx context.Context, id string, limit uint) ([]*TxLog, error) {
	if id == "" {
		return nil, fmt.Errorf("empty id")
	}
	if limit == 0 {
		return nil, fmt.Errorf("limit is 0")
	}

	s := n.d.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})

	results, err := s.Run(ctx, logQuery, map[string]any{"id": id}, neo4j.WithTxTimeout(txTimeout))
	if err != nil {
		return nil, fmt.Errorf("do query: %w", err)
	}

	logs, err := neo4j.CollectTWithContext(ctx, results, func(r *neo4j.Record) (*TxLog, error) {
		return txLog(r)
	})
	if err != nil {
		return nil, fmt.Errorf("collect logs: %w", err)
	}

	return logs, nil
}

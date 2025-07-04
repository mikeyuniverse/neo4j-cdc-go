package neo4j

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	cdcEarliestQuery = `CALL db.cdc.earliest()`
	logQuery         = `
		CALL db.cdc.query($id)
		YIELD id, txId, seq, metadata, event
		RETURN id, txId, seq, metadata, event
		LIMIT $limit
	`
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
	// Top-level fields
	ID       string         `json:"id"`
	TxID     int64          `json:"txId"`
	Seq      int64          `json:"seq"`
	Metadata *TxLogMetadata `json:"metadata,omitempty"`
	Event    *TxLogEvent    `json:"event,omitempty"`
}

// TxLogMetadata contains transaction metadata
type TxLogMetadata struct {
	DatabaseName      string
	ExecutingUser     string
	AuthenticatedUser string
	CaptureMode       string
	ConnectionClient  string
	ServerID          string
	ConnectionType    string
	ConnectionServer  string
	TxStartTime       time.Time
	TxCommitTime      time.Time
	TxMetadata        map[string]any
}

// TxLogEvent contains the actual change event data
type TxLogEvent struct {
	ElementID string                      `json:"elementId,omitempty"`
	Operation string                      `json:"operation,omitempty"` // "c" for create, "u" for update, "d" for delete
	Keys      map[string][]map[string]any `json:"keys,omitempty"`
	Labels    []string                    `json:"labels,omitempty"`    // For node events
	EventType string                      `json:"eventType,omitempty"` // "n" for node, "r" for relationship
	State     *TxLogState                 `json:"state,omitempty"`
}

// TxLogState contains before and after states
type TxLogState struct {
	Before *TxLogEntityState `json:"before,omitempty"`
	After  *TxLogEntityState `json:"after,omitempty"`
}

// TxLogEntityState represents the state of a node or relationship
type TxLogEntityState struct {
	Properties map[string]any `json:"properties,omitempty"`
	Labels     []string       `json:"labels,omitempty"` // For nodes
}

func (n *Neo4j) LogQuery(ctx context.Context, id string, limit uint) ([]*TxLog, error) {
	if id == "" {
		return nil, fmt.Errorf("empty id")
	}
	if limit == 0 {
		return nil, fmt.Errorf("limit is 0")
	}

	s := n.d.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})

	results, err := s.Run(ctx, logQuery, map[string]any{
		"id":    id,
		"limit": limit,
	}, neo4j.WithTxTimeout(txTimeout))
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

package neo4j

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	showDatabasesQuery = `
		SHOW DATABASES
		YIELD name
		WHERE name <> 'system'
		RETURN name`
)

func (n *Neo4j) DatabaseList(ctx context.Context) ([]string, error) {
	result, err := n.d.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"}).
		Run(ctx, showDatabasesQuery, nil, neo4j.WithTxTimeout(txTimeout))
	if err != nil {
		return nil, fmt.Errorf("do query: %w", err)
	}

	databases, err := neo4j.CollectTWithContext(ctx, result, func(r *neo4j.Record) (string, error) {
		return stringFromRecord(r, "name")
	})
	if err != nil {
		return nil, fmt.Errorf("collect results: %w", err)
	}

	return databases, nil
}

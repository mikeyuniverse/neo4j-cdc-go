package neo4j

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func stringFromRecord(r *neo4j.Record, key string) (string, error) {
	id, exists, err := neo4j.GetRecordValue[string](r, key)
	if err != nil {
		return "", fmt.Errorf("get %q: %w", key, err)
	}

	if !exists {
		return "", fmt.Errorf("%q is not found", key)
	}

	if id == "" {
		return "", fmt.Errorf("value is empty")
	}

	return id, nil
}

func txLog(r *neo4j.Record) (*TxLog, error) {
	return nil, nil
}

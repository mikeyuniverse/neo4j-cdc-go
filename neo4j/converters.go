package neo4j

import (
	"fmt"
	"time"

	"github.com/mikeyuniverse/neo4j-cdc-go/entities"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func stringFromRecord(r *neo4j.Record, key string) (string, error) {
	value, isNil, err := neo4j.GetRecordValue[string](r, key)
	if err != nil {
		return "", fmt.Errorf("get %q: %w", key, err)
	}

	if isNil {
		return "", fmt.Errorf("%q is not found", key)
	}

	if value == "" {
		return "", fmt.Errorf("value is empty")
	}

	return value, nil
}

func txLog(r *neo4j.Record) (*entities.TxLog, error) {
	log := &entities.TxLog{}

	// Parse top-level fields
	id, err := stringFromRecord(r, "id")
	if err != nil {
		return nil, fmt.Errorf("get id: %w", err)
	}
	log.ID = id

	txID, isNil, err := neo4j.GetRecordValue[int64](r, "txId")
	if err != nil {
		return nil, fmt.Errorf("get txId: %w", err)
	}
	if !isNil {
		log.TxID = txID
	}

	seq, isNil, err := neo4j.GetRecordValue[int64](r, "seq")
	if err != nil {
		return nil, fmt.Errorf("get seq: %w", err)
	}
	if !isNil {
		log.Seq = seq
	}

	// Parse metadata
	metadata, isNil, err := neo4j.GetRecordValue[map[string]any](r, "metadata")
	if err != nil {
		return nil, fmt.Errorf("get metadata: %w", err)
	}
	if !isNil {
		log.Metadata, err = parseTxLogMetadata(metadata)
		if err != nil {
			return nil, fmt.Errorf("parse metadata: %w", err)
		}
	}

	// Parse event
	event, isNil, err := neo4j.GetRecordValue[map[string]any](r, "event")
	if err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}
	if !isNil {
		log.Event, err = parseTxLogEvent(event)
		if err != nil {
			return nil, fmt.Errorf("parse event: %w", err)
		}
	}

	return log, nil
}

func parseTxLogMetadata(metadata map[string]any) (*entities.TxLogMetadata, error) {
	meta := &entities.TxLogMetadata{}

	if val, ok := metadata["executingUser"].(string); ok {
		meta.ExecutingUser = val
	}

	if val, ok := metadata["authenticatedUser"].(string); ok {
		meta.AuthenticatedUser = val
	}

	if val, ok := metadata["databaseName"].(string); ok {
		meta.DatabaseName = val
	}

	if val, ok := metadata["captureMode"].(string); ok {
		meta.CaptureMode = val
	}

	if val, ok := metadata["connectionClient"].(string); ok {
		meta.ConnectionClient = val
	}

	if val, ok := metadata["serverId"].(string); ok {
		meta.ServerID = val
	}

	if val, ok := metadata["connectionType"].(string); ok {
		meta.ConnectionType = val
	}

	if val, ok := metadata["connectionServer"].(string); ok {
		meta.ConnectionServer = val
	}

	if val, ok := metadata["txStartTime"].(time.Time); ok {
		meta.TxStartTime = val
	}

	if val, ok := metadata["txCommitTime"].(time.Time); ok {
		meta.TxCommitTime = val
	}

	if val, ok := metadata["txMetadata"].(map[string]any); ok {
		meta.TxMetadata = val
	}

	return meta, nil
}

func parseTxLogEvent(event map[string]any) (*entities.TxLogEvent, error) {
	evt := &entities.TxLogEvent{}

	if val, ok := event["elementId"].(string); ok {
		evt.ElementID = val
	}

	if val, ok := event["keys"].(map[string]any); ok {
		evt.Keys = make(map[string][]map[string]any)
		for k, v := range val {
			if keyList, ok := v.([]any); ok {
				keys := make([]map[string]any, len(keyList))
				for i, keyItem := range keyList {
					if keyMap, ok := keyItem.(map[string]any); ok {
						keys[i] = keyMap
					}
				}
				evt.Keys[k] = keys
			}
		}
	}

	if val, ok := event["eventType"].(string); ok {
		evt.EventType = entities.EventType(val)
	}

	if val, ok := event["state"].(map[string]any); ok {
		state, err := parseTxLogState(val)
		if err != nil {
			return nil, fmt.Errorf("parse state: %w", err)
		}
		evt.State = state
	}

	if val, ok := event["operation"].(string); ok {
		evt.Operation = entities.Operation(val)
	}

	if val, ok := event["labels"].([]any); ok {
		evt.Labels = make([]string, len(val))
		for i, label := range val {
			if labelStr, ok := label.(string); ok {
				evt.Labels[i] = labelStr
			}
		}
	}

	return evt, nil
}

func parseTxLogState(state map[string]any) (*entities.TxLogState, error) {
	s := &entities.TxLogState{}

	if val, ok := state["before"].(map[string]any); ok {
		before, err := parseTxLogEntityState(val)
		if err != nil {
			return nil, fmt.Errorf("parse before state: %w", err)
		}
		s.Before = before
	}

	if val, ok := state["after"].(map[string]any); ok {
		after, err := parseTxLogEntityState(val)
		if err != nil {
			return nil, fmt.Errorf("parse after state: %w", err)
		}
		s.After = after
	}

	return s, nil
}

func parseTxLogEntityState(entityState map[string]any) (*entities.TxLogEntityState, error) {
	state := &entities.TxLogEntityState{}

	if val, ok := entityState["properties"].(map[string]any); ok {
		state.Properties = val
	}

	if val, ok := entityState["labels"].([]any); ok {
		state.Labels = make([]string, len(val))
		for i, label := range val {
			if labelStr, ok := label.(string); ok {
				state.Labels[i] = labelStr
			}
		}
	}

	return state, nil
}

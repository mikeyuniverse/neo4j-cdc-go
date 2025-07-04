package entities

import "time"

type EventType string

func (e EventType) String() string {
	switch e {
	case Node:
		return "Node"
	case Relation:
		return "Relation"
	default:
		return string(e)
	}
}

type Operation string

func (o Operation) String() string {
	switch o {
	case Created:
		return "Created"
	case Updated:
		return "Updated"
	case Deleted:
		return "Deleted"
	default:
		return string(o)
	}
}

const (
	Node     EventType = "n"
	Relation EventType = "r"

	Created Operation = "c"
	Updated Operation = "u"
	Deleted Operation = "d"
)

type TxLog struct {
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
	Operation Operation                   `json:"operation,omitempty"`
	EventType EventType                   `json:"eventType,omitempty"`
	Keys      map[string][]map[string]any `json:"keys,omitempty"`
	Labels    []string                    `json:"labels,omitempty"`
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
	Labels     []string       `json:"labels,omitempty"`
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/mikeyuniverse/neo4j-cdc-go/neo4j"
)

func main() {
	slog.Info("starting application...")

	ctx := context.Background()

	n := neo4j.New()
	err := n.Connect(ctx)
	if err != nil {
		exitErr(err)
		return
	}

	earliestID, err := n.Earliest(ctx)
	if err != nil {
		exitErr(err)
		return
	}

	logs, err := n.LogQuery(ctx, earliestID, 2)
	if err != nil {
		exitErr(err)
		return
	}

	for _, l := range logs {
		fmt.Printf("%s - %d - Type: %s | Operation: %s\n", l.ID, l.TxID, l.Event.EventType, l.Event.Operation)
	}
}

type LogType int

const (
	NodeCreated = iota + 1
	NodeUpdated
	NodeDeleted
	RelCreated
	RelUpdated
	RelDeleted
)

func Switcher(ctx context.Context, log *entities.TxLog) error {
	if log.Event == nil {
		return fmt.Errorf("event is nil for log with id %q and metadata %+v", log.ID, log.Metadata)
	}
	
	if log.Event.EventType == entities.Node {
		
	} else if log.Event.EventType == entities.Relation

	switch log.Event.EventType {
	case entities.Node:
		return nil
	case entities.Relation:
		return nil
	default:
		return fmt.Errorf("unknown event type: %+v", log)
	}
}

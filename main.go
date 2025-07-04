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

func exitErr(err error) {
	slog.Error(err.Error())
	os.Exit(1)
}

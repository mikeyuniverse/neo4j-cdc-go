package main

import (
	"context"
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

	_, err = n.Earliest(ctx)
	if err != nil {
		exitErr(err)
		return
	}

}

func exitErr(err error) {
	slog.Error(err.Error())
	os.Exit(1)
}

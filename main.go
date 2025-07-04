package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mikeyuniverse/neo4j-cdc-go/nats"
	"github.com/mikeyuniverse/neo4j-cdc-go/neo4j"
	"github.com/mikeyuniverse/neo4j-cdc-go/proc"
)

func main() {
	setSlogLevel()

	slog.Info("starting application...")

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGKILL)
	defer cancel()

	n := neo4j.New()
	err := n.Connect(ctx)
	if err != nil {
		slog.Error("on neo4j connecting", "error", err.Error())
		return
	}
	defer n.Close(ctx) // TODO: check for already cancelled context

	q := nats.New()
	err = q.Connect(ctx)
	if err != nil {
		slog.Error("on nats connecting", "error", err.Error())
		return
	}
	defer q.Close(ctx) // TODO: check for already cancelled context

	publishAll := PublishAllHandler{
		subject: "cdc.event",
		p:       q,
	}

	p := proc.New(
		ctx, n, q,
		proc.TimerMiddleware(
			proc.LoggerMiddleware(
				proc.Matcher(publishAll),
			),
		),
	)

	err = p.Run(ctx)
	if err != nil {
		slog.Error("error on run", "error", err)
		return
	}

	slog.Info("application finished without errors")
}

func setSlogLevel() {
	lvl := os.Getenv("SLOG_LEVEL")

	data := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	slog.SetLogLoggerLevel(data[strings.ToLower(lvl)])
}

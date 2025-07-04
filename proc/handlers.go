package proc

import (
	"context"
	"log/slog"
	"time"

	"github.com/mikeyuniverse/neo4j-cdc-go/entities"
)

type Handler func(ctx context.Context, log *entities.TxLog) error

type Middleware func(next Handler) Handler

func LoggerMiddleware(next Handler) Handler {
	return func(ctx context.Context, log *entities.TxLog) error {
		printLog(log)
		return next(ctx, log)
	}
}

func TimerMiddleware(next Handler) Handler {
	return func(ctx context.Context, log *entities.TxLog) error {
		start := time.Now()
		defer func() {
			slog.Info("timer", "id", log.ID, "duration", time.Since(start).String())
		}()

		return next(ctx, log)
	}
}

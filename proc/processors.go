package proc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mikeyuniverse/neo4j-cdc-go/entities"
)

type Processor interface {
	OnNodeCreated(ctx context.Context, log *entities.TxLog) error
	OnNodeUpdated(ctx context.Context, log *entities.TxLog) error
	OnNodeDeleted(ctx context.Context, log *entities.TxLog) error

	OnRelationCreated(ctx context.Context, log *entities.TxLog) error
	OnRelationUpdated(ctx context.Context, log *entities.TxLog) error
	OnRelationDeleted(ctx context.Context, log *entities.TxLog) error
}

func Matcher(p Processor) Handler {
	return func(ctx context.Context, log *entities.TxLog) error {
		return match(ctx, log, p)
	}
}

// FIXME: replace matching paths by input log
func match(ctx context.Context, log *entities.TxLog, p Processor) error {
	if p == nil {
		return errors.New("processor not defined")
	}

	if log.Event == nil {
		return fmt.Errorf("event is nil for log with id %q and metadata %+v", log.ID, log.Metadata)
	}

	paths := map[string]Handler{
		string(entities.Created) + string(entities.Node):     p.OnNodeCreated,
		string(entities.Updated) + string(entities.Node):     p.OnNodeUpdated,
		string(entities.Deleted) + string(entities.Node):     p.OnNodeDeleted,
		string(entities.Created) + string(entities.Relation): p.OnRelationCreated,
		string(entities.Updated) + string(entities.Relation): p.OnRelationUpdated,
		string(entities.Deleted) + string(entities.Relation): p.OnRelationDeleted,
	}

	h, found := paths[string(log.Event.Operation)+string(log.Event.EventType)]
	if !found {
		return fmt.Errorf("handler not found for log %q - %q - %q", log.Event.EventType, log.Event.Operation, log.ID)
	}

	return h(ctx, log)
}

// start N workers to handle logs
func (p *Proc) goProcessing(workersCount int) {
	for range workersCount {
		p.pool.Go(func(c context.Context) error {
			return p.handleLogs(c, p.txs)
		})
	}
}

func (p *Proc) handleLogs(ctx context.Context, src <-chan *entities.TxLog) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case log, ok := <-src:
			if !ok {
				return nil
			}

			p.h(ctx, log)
		}
	}
}

func printLog(log *entities.TxLog) {
	// data, err := json.MarshalIndent(log, "", "  ")
	data, err := json.Marshal(log)
	if err != nil {
		slog.Error("marshal error", "error", err)
		return
	}

	slog.Info(string(data))
}

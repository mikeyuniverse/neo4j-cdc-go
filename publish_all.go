package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mikeyuniverse/neo4j-cdc-go/entities"
)

type Publisher interface {
	Publish(ctx context.Context, id string, subject string, payload []byte) error
}

type PublishAllHandler struct {
	subject string
	p       Publisher
}

func (n PublishAllHandler) OnNodeCreated(ctx context.Context, log *entities.TxLog) error {
	return n.publish(ctx, log)
}

func (n PublishAllHandler) OnNodeUpdated(ctx context.Context, log *entities.TxLog) error {
	return n.publish(ctx, log)
}

func (n PublishAllHandler) OnNodeDeleted(ctx context.Context, log *entities.TxLog) error {
	return n.publish(ctx, log)
}

func (n PublishAllHandler) OnRelationCreated(ctx context.Context, log *entities.TxLog) error {
	return n.publish(ctx, log)
}

func (n PublishAllHandler) OnRelationUpdated(ctx context.Context, log *entities.TxLog) error {
	return n.publish(ctx, log)
}

func (n PublishAllHandler) OnRelationDeleted(ctx context.Context, log *entities.TxLog) error {
	return n.publish(ctx, log)
}

func (n PublishAllHandler) publish(ctx context.Context, log *entities.TxLog) error {
	data, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("marshal tx log: %w", err)
	}

	return n.p.Publish(ctx, log.ID, n.subject, data)
}

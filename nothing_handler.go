package main

import (
	"context"

	"github.com/mikeyuniverse/neo4j-cdc-go/entities"
)

type NothingHandler struct{}

func (n NothingHandler) OnNodeCreated(ctx context.Context, log *entities.TxLog) error     { return nil }
func (n NothingHandler) OnNodeUpdated(ctx context.Context, log *entities.TxLog) error     { return nil }
func (n NothingHandler) OnNodeDeleted(ctx context.Context, log *entities.TxLog) error     { return nil }
func (n NothingHandler) OnRelationCreated(ctx context.Context, log *entities.TxLog) error { return nil }
func (n NothingHandler) OnRelationUpdated(ctx context.Context, log *entities.TxLog) error { return nil }
func (n NothingHandler) OnRelationDeleted(ctx context.Context, log *entities.TxLog) error { return nil }

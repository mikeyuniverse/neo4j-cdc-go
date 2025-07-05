package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/mikeyuniverse/neo4j-cdc-go/entities"
)

type CounterHandler struct {
	mu   sync.Mutex
	data map[string]int64
}

func (n *CounterHandler) OnNodeCreated(ctx context.Context, log *entities.TxLog) error {
	n.inc(log.Metadata.DatabaseName)
	return nil
}

func (n *CounterHandler) OnNodeUpdated(ctx context.Context, log *entities.TxLog) error {
	n.inc(log.Metadata.DatabaseName)
	return nil
}

func (n *CounterHandler) OnNodeDeleted(ctx context.Context, log *entities.TxLog) error {
	n.inc(log.Metadata.DatabaseName)
	return nil
}

func (n *CounterHandler) OnRelationCreated(ctx context.Context, log *entities.TxLog) error {
	n.inc(log.Metadata.DatabaseName)
	return nil
}

func (n *CounterHandler) OnRelationUpdated(ctx context.Context, log *entities.TxLog) error {
	n.inc(log.Metadata.DatabaseName)
	return nil
}

func (n *CounterHandler) OnRelationDeleted(ctx context.Context, log *entities.TxLog) error {
	n.inc(log.Metadata.DatabaseName)
	return nil
}

func (n *CounterHandler) inc(name string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	value := n.data[name]
	n.data[name] = value + 1
}

func (n *CounterHandler) Print() {
	n.mu.Lock()
	defer n.mu.Unlock()

	for k, v := range n.data {
		fmt.Println(" - ", k, " ::", v)
	}
}

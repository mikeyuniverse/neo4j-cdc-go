package proc

import (
	"context"
	"fmt"

	"github.com/mikeyuniverse/neo4j-cdc-go/entities"
	"github.com/sourcegraph/conc/pool"
)

const (
	logPageSize = 10
	bufferSize  = 100
	workers     = 10
)

type Graph interface {
	DatabaseList(ctx context.Context) ([]string, error)
	EnableFullCDC(ctx context.Context, databases []string) error

	Earliest(ctx context.Context, database string) (string, error)
	GetCDCItem(ctx context.Context, database, startAtID string, logPageSize uint) ([]*entities.TxLog, error)
}

type Publisher interface {
	Publish(ctx context.Context, id, subject string, payload []byte) error
}

type Proc struct {
	graph Graph
	q     Publisher
	pool  *pool.ContextPool
	txs   chan *entities.TxLog
	h     Handler
}

// TODO: set goroutines limits
// pool.WithMaxGoroutines(workers + len(databases))
func New(ctx context.Context, graph Graph, q Publisher, h Handler) *Proc {
	pool := pool.New().
		WithContext(ctx).
		WithCancelOnError().
		WithFirstError()

	return &Proc{
		graph: graph,
		q:     q,
		pool:  pool,
		txs:   make(chan *entities.TxLog, bufferSize),
		h:     h,
	}
}

func (p *Proc) Run(ctx context.Context) error {
	// TODO: remove this functions from this place
	dbs, err := p.graph.DatabaseList(ctx)
	if err != nil {
		return fmt.Errorf("get databases: %w", err)
	}

	err = p.graph.EnableFullCDC(ctx, dbs)
	if err != nil {
		return fmt.Errorf("enable CDC: %w", err)
	}

	return p.startCDCProcessing(ctx, dbs)
}

func (p *Proc) startCDCProcessing(ctx context.Context, databases []string) error {
	p.goReading(ctx, databases)
	p.goProcessing(workers)

	err := p.pool.Wait()
	if err != nil {
		return fmt.Errorf("wait for goroutines: %w", err)
	}

	return nil
}

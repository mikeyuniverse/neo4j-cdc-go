package proc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

// FIXME: there can be more than one (or two) database
// we need to determine how many goroutines can be started
// for reading transactions log
//
// start goroutine per database
func (p *Proc) goReading(ctx context.Context, databases []string) {
	for _, db := range databases {
		p.pool.Go(func(c context.Context) error {
			earliestID, err := p.graph.Earliest(ctx, db)
			if err != nil {
				return fmt.Errorf("get earliest ID for %q: %w", db, err)
			}

			slog.Info(
				"start reading log",
				"database", db,
				"start", earliestID,
			)

			// blocking call
			err = p.readQueries(c, db, earliestID)
			return err
		})
	}
}

func (p *Proc) readQueries(ctx context.Context, db, logID string) error {
	tickInterval := time.Second

	t := time.NewTicker(tickInterval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
		}

		logs, err := p.graph.GetCDCItem(ctx, db, logID, logPageSize)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				continue
			}
			return fmt.Errorf("do log query: %w", err)
		}

		if len(logs) > 0 {
			logID = logs[len(logs)-1].ID
		}

		for _, l := range logs {
			p.txs <- l
		}

		t.Reset(tickInterval)
	}
}

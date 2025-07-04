package nats

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
	js "github.com/nats-io/nats.go/jetstream"
)

type Nats struct {
	connect *nats.Conn
	js      js.JetStream
}

func New() *Nats {
	return &Nats{}
}

func (n *Nats) Connect(ctx context.Context) error {
	connectionURL := nats.DefaultURL
	if v := os.Getenv("NATS_DSN"); v != "" {
		connectionURL = v
	}

	c, err := nats.Connect(connectionURL)
	if err != nil {
		return fmt.Errorf("connect %q: %w", connectionURL, err)
	}

	js, err := js.New(c)
	if err != nil {
		return fmt.Errorf("jet stream: %w", err)
	}

	n.connect = c
	n.js = js

	return nil
}

func (n *Nats) Close(ctx context.Context) error {
	defer n.connect.Close()
	err := n.connect.FlushWithContext(ctx)
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			return err
		}
	}
	return nil
}

func (n *Nats) Publish(ctx context.Context, id, subject string, payload []byte) error {
	_, err := n.js.Publish(ctx, subject, payload, js.WithMsgID(id))
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

package neo4j

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	txTimeout = time.Second * 5
)

var (
	connectDSN = os.Getenv("NEO4J_DSN")
	username   = os.Getenv("NEO4J_USERNAME")
	password   = os.Getenv("NEO4J_PASSWORD")
)

type Neo4j struct {
	d neo4j.DriverWithContext
}

func New() *Neo4j {
	return &Neo4j{}
}

func (n *Neo4j) Connect(ctx context.Context) error {
	driver, err := neo4j.NewDriverWithContext(connectDSN, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return fmt.Errorf("new driver: %w", err)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return fmt.Errorf("verify connectivity: %w", err)
	}

	n.d = driver

	return nil
}

func (n *Neo4j) Close(ctx context.Context) error {
	return n.d.Close(ctx)
}

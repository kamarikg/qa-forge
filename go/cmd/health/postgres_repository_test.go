package main

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestPostgreSQLRepositoryConnects(t *testing.T) {
	conn, err := pgx.Connect(
		context.Background(),
		"postgres://qa_forge:qa_forge@localhost:5432/qa_forge",
	)
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close(context.Background())

	if err := conn.Ping(context.Background()); err != nil {
		t.Fatalf("failed to ping PostgreSQL: %v", err)
	}
}

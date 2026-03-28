package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		fmt.Println("TEST_DATABASE_URL not set, skipping integration tests")
		os.Exit(0)
	}

	var err error
	pool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		fmt.Printf("failed to create pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := runMigrations(pool); err != nil {
		fmt.Printf("failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func runMigrations(pool *pgxpool.Pool) error {
	//TODO в production ready использовать мигратор и через него будем доводить бд до нужного состояния
	_, err := pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS invitations (
			id         BIGSERIAL    PRIMARY KEY,
			code       VARCHAR(255) NOT NULL UNIQUE,
			max_uses   INT          NOT NULL,
			used_count INT          NOT NULL
		);

		CREATE TABLE IF NOT EXISTS registrations (
			id            BIGSERIAL    PRIMARY KEY,
			invitation_id BIGINT       NOT NULL REFERENCES invitations(id),
			email         VARCHAR(255) NOT NULL,
			created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			UNIQUE (invitation_id, email)
		);
	`)
	return err
}

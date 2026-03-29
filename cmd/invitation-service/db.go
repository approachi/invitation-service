package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initDB(dbURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	log.Println("connected to database")
	return pool, nil
}

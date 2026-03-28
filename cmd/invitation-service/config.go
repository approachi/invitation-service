package main

import (
	"fmt"
	"os"
)

type config struct {
	DatabaseURL string
	Addr        string
}

func loadConfig() (config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return config{}, fmt.Errorf("DATABASE_URL is required")
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":3000"
	}

	return config{
		DatabaseURL: dbURL,
		Addr:        addr,
	}, nil
}

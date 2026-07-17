// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/andreychh/coopera/internal/api"
	"github.com/andreychh/coopera/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL())
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	mux := http.NewServeMux()
	mux.HandleFunc(
		"GET /health",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	strict := api.NewStrictHandler(api.NewServer(domain.NewSQLWorld(pool)), nil)
	handler := api.HandlerFromMuxWithBaseURL(strict, mux, "/v1")

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return server.ListenAndServe()
}

func databaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
}

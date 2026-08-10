// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/andreychh/coopera/internal/api"
	"github.com/andreychh/coopera/internal/post"
	"github.com/andreychh/coopera/internal/seal"
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

	key, err := accessTokenKey()
	if err != nil {
		return err
	}

	strict := api.NewStrictHandlerWithOptions(
		api.NewServer(pool, post.NewLog(), seal.NewJWT(key)),
		nil,
		api.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  api.RequestError,
			ResponseErrorHandlerFunc: api.ResponseError,
		},
	)
	handler := api.HandlerWithOptions(strict, api.StdHTTPServerOptions{
		BaseURL:     "/v1",
		BaseRouter:  mux,
		Middlewares: []api.MiddlewareFunc{api.NewGate(pool)},
	})

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

// accessTokenKey reads the key access tokens are signed with. There is
// no default and cannot be: a key shipped in the source would be known
// to everyone who has read it, and anybody could then write themselves a
// pass as anybody else. Refusing to start is the only safe way to be
// missing it.
//
// Thirty-two bytes is the floor because the signature is HMAC-SHA256,
// which is no stronger than the key behind it; a short one can be
// searched for offline against any single token that has ever been
// issued.
func accessTokenKey() ([]byte, error) {
	key, exists := os.LookupEnv("ACCESS_TOKEN_KEY")
	if !exists {
		return nil, errors.New("ACCESS_TOKEN_KEY is not set")
	}
	if len(key) < 32 {
		return nil, errors.New("ACCESS_TOKEN_KEY must be at least 32 bytes")
	}
	return []byte(key), nil
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

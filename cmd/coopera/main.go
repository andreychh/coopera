// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/netip"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andreychh/coopera/internal/api/v2/identified"
	"github.com/andreychh/coopera/internal/api/v2/limit"
	"github.com/andreychh/coopera/internal/api/v2/located"
	"github.com/andreychh/coopera/internal/api/v2/raw"
	"github.com/andreychh/coopera/internal/ratelimiting"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := run(ctx)
	stop()

	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	network, err := ratelimiting.NewPolicy(20, 100*time.Millisecond)
	if err != nil {
		return fmt.Errorf("network rate: %w", err)
	}
	networks, err := ratelimiting.NewBuckets[netip.Prefix](network, 10_000)
	if err != nil {
		return fmt.Errorf("network buckets: %w", err)
	}
	overall, err := ratelimiting.NewPolicy(1000, time.Millisecond)
	if err != nil {
		return fmt.Errorf("service rate: %w", err)
	}
	service := ratelimiting.NewBucket(overall)

	const (
		// body is what a JSON route accepts, and timeout is how long it may
		// take. A route that takes uploads or runs long passes its own bounds
		// to raw.NewCheckpoint.
		body    = 1 << 20
		timeout = 10 * time.Second
		// grace is how long a request that is already being served may take to
		// finish once the server is told to stop.
		grace = 20 * time.Second
	)

	// stages wraps a handler in every stage and filter a request passes
	// before it. The outer limit runs first, so a request its network refuses
	// takes nothing from the service. All routes share the buckets.
	stages := func(next located.Handler) http.Handler {
		return raw.NewCheckpoint(
			identified.NewCheckpoint(
				located.NewCheckpoint(
					limit.Keyed(
						limit.Shared(next, service, http.StatusServiceUnavailable),
						networks,
						located.Request.ClientNetwork,
						http.StatusTooManyRequests,
					),
				),
			),
			body,
			timeout,
		)
	}

	health := located.HandlerFunc(
		func(w http.ResponseWriter, r located.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	mux := http.NewServeMux()
	mux.Handle("GET /health", stages(health))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		// Headers are small and arrive at once, so a client that takes longer
		// is stalling the connection on purpose.
		ReadHeaderTimeout: 5 * time.Second,
		// A body of the size the API accepts takes this long to arrive over a
		// slow mobile link.
		ReadTimeout: 15 * time.Second,
		// Longer than the deadline a request is given, so the server answers
		// before it cuts the connection.
		WriteTimeout: 15 * time.Second,
		// Longer than the idle timeout of a proxy in front, which would
		// otherwise reuse a connection just as this server closes it.
		IdleTimeout: 120 * time.Second,
		// Above the largest headers a route takes, a few kilobytes of access
		// token among them. It is read before routing, so it is one number for
		// every route.
		MaxHeaderBytes: 64 << 10,

		MaxHeaderValueCount: 500,
	}

	return serve(ctx, server, grace)
}

// serve runs server until ctx is done, and then gives the requests it is
// already serving grace to finish.
func serve(ctx context.Context, server *http.Server, grace time.Duration) error {
	errs := make(chan error, 1)
	go func() { errs <- server.ListenAndServe() }()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
	}

	// The signal has already canceled ctx, so the shutdown is given a deadline
	// of its own.
	stopping, cancel := context.WithTimeout(context.WithoutCancel(ctx), grace)
	defer cancel()
	return server.Shutdown(stopping)
}

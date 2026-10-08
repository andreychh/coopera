// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package located_test

import (
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"
	"time"

	"github.com/andreychh/coopera/internal/api/v2/identified"
	"github.com/andreychh/coopera/internal/api/v2/located"
	"github.com/andreychh/coopera/internal/api/v2/raw"
)

func TestCheckpoint_ServeHTTP(t *testing.T) {
	t.Parallel()

	// serve passes a request that arrived from remote through the checkpoint
	// and returns the response status, whether the handler was reached and the
	// client address it received.
	serve := func(t *testing.T, remote string) (int, bool, netip.Addr) {
		t.Helper()

		var (
			reached bool
			client  netip.Addr
		)
		checkpoint := raw.NewCheckpoint(identified.NewCheckpoint(located.NewCheckpoint(
			located.HandlerFunc(func(_ http.ResponseWriter, r located.Request) {
				reached, client = true, r.ClientAddr()
			}),
		)), 1<<20, time.Minute)
		request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", http.NoBody)
		request.RemoteAddr = remote
		recorder := httptest.NewRecorder()
		checkpoint.ServeHTTP(recorder, request)

		return recorder.Code, reached, client
	}

	passed := []struct {
		behavior string
		remote   string
		client   netip.Addr
	}{{
		behavior: "passes on an IPv4 address",
		remote:   "203.0.113.7:51234",
		client:   netip.MustParseAddr("203.0.113.7"),
	}, {
		behavior: "passes on an IPv6 address",
		remote:   "[2001:db8::7]:51234",
		client:   netip.MustParseAddr("2001:db8::7"),
	}, {
		behavior: "passes on an IPv4-mapped address as IPv4",
		remote:   "[::ffff:203.0.113.7]:51234",
		client:   netip.MustParseAddr("203.0.113.7"),
	}}
	for _, c := range passed {
		t.Run(c.behavior, func(t *testing.T) {
			t.Parallel()

			status, reached, client := serve(t, c.remote)

			if status != http.StatusOK || !reached || client != c.client {
				t.Errorf(
					"from %q: status %d, reached %t, ClientAddr() = %v, want 200, reached, %v",
					c.remote, status, reached, client, c.client,
				)
			}
		})
	}

	refused := []struct {
		behavior string
		remote   string
	}{{
		behavior: "answers 500 to an address without a port",
		remote:   "203.0.113.7",
	}, {
		behavior: "answers 500 to an empty address, as on a Unix socket",
		remote:   "",
	}}
	for _, c := range refused {
		t.Run(c.behavior, func(t *testing.T) {
			t.Parallel()

			status, reached, _ := serve(t, c.remote)

			if status != http.StatusInternalServerError || reached {
				t.Errorf(
					"from %q: status %d, reached %t, want 500, not reached",
					c.remote, status, reached,
				)
			}
		})
	}
}

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package raw_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andreychh/coopera/internal/api/v2/raw"
)

func TestRequest_Context(t *testing.T) {
	t.Parallel()

	t.Run("carries the values of the original context", func(t *testing.T) {
		t.Parallel()

		type key struct{}
		ctx := context.WithValue(t.Context(), key{}, "carried")
		origin := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", http.NoBody)

		got := serve(t, "/", origin).Context().Value(key{})

		if got != "carried" {
			t.Errorf("Context().Value = %v, want the value of the original context", got)
		}
	})
}

func TestRequest_Header(t *testing.T) {
	t.Parallel()

	t.Run("matches the name without regard to case", func(t *testing.T) {
		t.Parallel()

		const token = "Bearer token"
		origin := newRequest(t, "/")
		origin.Header.Set("Authorization", token)

		got := serve(t, "/", origin).Header("authorization")

		if got != token {
			t.Errorf("Header(authorization) = %q, want %q", got, token)
		}
	})

	t.Run("returns nothing for a header the request lacks", func(t *testing.T) {
		t.Parallel()

		got := serve(t, "/", newRequest(t, "/")).Header("Authorization")

		if got != "" {
			t.Errorf("Header(Authorization) = %q, want none", got)
		}
	})
}

func TestRequest_PathValue(t *testing.T) {
	t.Parallel()

	t.Run("returns the value of a wildcard in the pattern", func(t *testing.T) {
		t.Parallel()

		got := serve(t, "GET /teams/{id}", newRequest(t, "/teams/42")).PathValue("id")

		if got != "42" {
			t.Errorf("PathValue(id) = %q, want %q", got, "42")
		}
	})

	t.Run("returns nothing for a wildcard the pattern lacks", func(t *testing.T) {
		t.Parallel()

		got := serve(t, "GET /teams/{id}", newRequest(t, "/teams/42")).PathValue("name")

		if got != "" {
			t.Errorf("PathValue(name) = %q, want none", got)
		}
	})
}

func TestRequest_Query(t *testing.T) {
	t.Parallel()

	cases := []struct {
		behavior string
		name     string
		value    string
		present  bool
	}{{
		behavior: "returns the value of a parameter the request carries",
		name:     "status",
		value:    "active",
		present:  true,
	}, {
		behavior: "reports a parameter carried without a value as present",
		name:     "blank",
		value:    "",
		present:  true,
	}, {
		behavior: "reports a parameter the request lacks as absent",
		name:     "absent",
		value:    "",
		present:  false,
	}}
	for _, c := range cases {
		t.Run(c.behavior, func(t *testing.T) {
			t.Parallel()

			origin := newRequest(t, "/links?status=active&blank=")

			value, present := serve(t, "/", origin).Query(c.name)

			if value != c.value || present != c.present {
				t.Errorf(
					"Query(%s) = (%q, %t), want (%q, %t)",
					c.name, value, present, c.value, c.present,
				)
			}
		})
	}
}

func TestRequest_Body(t *testing.T) {
	t.Parallel()

	t.Run("yields what the client sent", func(t *testing.T) {
		t.Parallel()

		origin := httptest.NewRequestWithContext(
			t.Context(),
			http.MethodPost,
			"/",
			strings.NewReader("payload"),
		)

		got, err := io.ReadAll(serve(t, "/", origin).Body())
		if err != nil {
			t.Fatalf("read Body: %v", err)
		}
		if string(got) != "payload" {
			t.Errorf("Body = %q, want %q", got, "payload")
		}
	})
}

func newRequest(t *testing.T, target string) *http.Request {
	t.Helper()

	return httptest.NewRequestWithContext(t.Context(), http.MethodGet, target, http.NoBody)
}

// serve routes origin through a Checkpoint registered for pattern and returns
// the request its handler received. It fails the test if origin does not match
// pattern.
func serve(t *testing.T, pattern string, origin *http.Request) raw.Request {
	t.Helper()

	var (
		got    raw.Request
		served bool
	)
	mux := http.NewServeMux()
	mux.Handle(pattern, raw.NewCheckpoint(raw.HandlerFunc(
		func(_ http.ResponseWriter, r raw.Request) {
			got, served = r, true
		},
	)))
	mux.ServeHTTP(httptest.NewRecorder(), origin)

	if !served {
		t.Fatalf("%s %s did not reach the handler", origin.Method, origin.URL)
	}
	return got
}

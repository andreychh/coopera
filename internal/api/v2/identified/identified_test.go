// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package identified_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreychh/coopera/internal/api/v2/identified"
	"github.com/andreychh/coopera/internal/api/v2/raw"
	"github.com/google/uuid"
)

func TestCheckpoint_ServeHTTP(t *testing.T) {
	t.Parallel()

	// serve passes request through the checkpoint and returns the identifier
	// the handler received and the one reported in the X-Request-ID header.
	serve := func(request *http.Request) (uuid.UUID, string) {
		var id uuid.UUID
		checkpoint := raw.NewCheckpoint(identified.NewCheckpoint(identified.HandlerFunc(
			func(_ http.ResponseWriter, r identified.Request) {
				id = r.ID()
			},
		)))
		recorder := httptest.NewRecorder()
		checkpoint.ServeHTTP(recorder, request)

		return id, recorder.Header().Get("X-Request-ID")
	}
	newRequest := func(t *testing.T) *http.Request {
		t.Helper()

		return httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", http.NoBody)
	}

	t.Run("reports the identifier it gives", func(t *testing.T) {
		t.Parallel()

		id, header := serve(newRequest(t))

		if id == uuid.Nil {
			t.Fatal("ID() is nil, want a generated identifier")
		}
		if header != id.String() {
			t.Errorf("X-Request-ID = %q, want %q", header, id)
		}
	})

	t.Run("ignores an identifier the client chose", func(t *testing.T) {
		t.Parallel()

		chosen := uuid.New().String()
		request := newRequest(t)
		request.Header.Set("X-Request-ID", chosen)

		id, header := serve(request)

		if id.String() == chosen {
			t.Errorf("ID() = %s, want one the server generated", id)
		}
		if header == chosen {
			t.Errorf("X-Request-ID = %q, want the server's identifier", header)
		}
	})

	t.Run("gives each request its own identifier", func(t *testing.T) {
		t.Parallel()

		first, _ := serve(newRequest(t))
		second, _ := serve(newRequest(t))

		if first == second {
			t.Errorf("two requests share the identifier %s", first)
		}
	})
}

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewGate builds the middleware that turns away callers who have not
// introduced themselves.
//
// Four operations stay open to them: looking at themselves, giving
// themselves a name, renewing their pass and leaving. Everything else
// shows a person to other people, and there is nothing yet to show.
//
// The last two are open out of necessity rather than kindness. Without
// renewal a pass would die fifteen minutes in and take with it the only
// thing its holder came to do; and leaving is nobody's to withhold —
// a person who wants out gets out, named or not.
//
// Matching is by exact path because all four are fixed, without
// parameters; anything unlisted is closed, so a new operation is guarded
// by default and forgetting about it costs a refusal rather than a hole.
//
// The gate answers outside the typed responses the generator builds per
// operation: those live inside the strict handler, and this stands in
// front of it. The body is the same either way — RFC 9457 problem
// details.
//
// A caller whose X-User-Id cannot be read passes through untouched. The
// gate has nobody to ask about, and the handler behind it answers 400
// for that same reason, so nothing gets in that would not have anyway.
func NewGate(pool *pgxpool.Pool) MiddlewareFunc {
	open := map[string]string{
		"/v1/users/me":                      http.MethodGet,
		"/v1/users/me/introduction":         http.MethodPost,
		"/v1/auth/sessions/current":         http.MethodDelete,
		"/v1/auth/sessions/current/refresh": http.MethodPost,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if method, listed := open[r.URL.Path]; listed && method == r.Method {
				next.ServeHTTP(w, r)
				return
			}

			actorID, err := domain.ParseID(r.Header.Get("X-User-Id"))
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			introduced, checkErr := usecase.NewIsIntroduced(pool, actorID).Exec(r.Context())
			if checkErr != nil {
				writeProblem(w, gateFailure(checkErr), "")
				return
			}
			if !introduced {
				writeProblem(w, http.StatusForbidden, "Introduce yourself first")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// gateFailure is a type switch rather than a chain of checks so that
// gochecksumtype verifies it: a failure added to
// [domain.IsIntroducedError] breaks the build here instead of quietly
// becoming a 500.
func gateFailure(err domain.IsIntroducedError) int {
	switch err.(type) {
	case domain.UserNotFoundError:
		return http.StatusUnauthorized
	case domain.UnexpectedError:
		return http.StatusInternalServerError
	}
	return http.StatusInternalServerError
}

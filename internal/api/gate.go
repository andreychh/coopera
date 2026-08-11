// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

// gate is the doorkeeper every request passes on its way to a handler.
//
// Doors come in three kinds, and the two maps below name the first two;
// everything unlisted is of the third.
//
// Three stand open to anyone. Asking for a code and signing in are how a
// person comes by a pass at all, so demanding one there would shut the
// only way in. Renewing is open for a sharper reason: the access token
// is expired at exactly the moment renewal is wanted, and what is shown
// instead travels in the body and is judged by the usecase itself.
//
// Three ask for a pass but not for a name. Looking at oneself and naming
// oneself are how a person stops being nameless, and leaving is nobody's
// to withhold — someone who wants out gets out, named or not.
//
// Everything else asks for both. This is not a punishment: almost every
// other act shows a person to somebody else, and there is nothing to
// show yet.
//
// Anything unlisted falls to the strictest kind, so a new operation is
// guarded by default and forgetting about it costs a refusal rather than
// a hole.
type gate struct {
	pool *pgxpool.Pool
	seal Seal

	unguarded map[string]string
	nameless  map[string]string
}

func NewGate(pool *pgxpool.Pool, seal Seal) MiddlewareFunc {
	g := gate{
		pool: pool,
		seal: seal,
		unguarded: map[string]string{
			"/v1/auth/codes":                    http.MethodPost,
			"/v1/auth/sessions":                 http.MethodPost,
			"/v1/auth/sessions/current/refresh": http.MethodPost,
		},
		nameless: map[string]string{
			"/v1/users/me":              http.MethodGet,
			"/v1/users/me/introduction": http.MethodPost,
			"/v1/auth/sessions/current": http.MethodDelete,
		},
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if matches(g.unguarded, r) {
				next.ServeHTTP(w, r)
				return
			}

			actor, admitted := g.admit(w, r)
			if !admitted {
				return
			}

			next.ServeHTTP(w, r.WithContext(withActor(r.Context(), actor)))
		})
	}
}

// admit says who is at the door, or turns them away itself. A false
// second result means the answer has already been written and there is
// nothing further to do with the request.
//
// The pass is looked up on every guarded request rather than believed
// for the fifteen minutes it claims. A signature cannot be taken back,
// so signing out and catching a theft would otherwise take effect only
// once the token died of old age. One question answers both what the
// session is worth now and whether its holder has a name.
//
// The answers are given outside the typed responses the generator builds
// per operation: those live inside the strict handler, and this stands
// in front of it. The body is the same either way — RFC 9457 problem
// details.
func (g gate) admit(w http.ResponseWriter, r *http.Request) (domain.Actor, bool) {
	token, shown := bearerToken(r)
	if !shown {
		writeProblem(w, http.StatusUnauthorized, "Show a pass")
		return domain.Actor{}, false
	}

	actor, err := g.seal.Read(token)
	if err != nil {
		writeProblem(w, http.StatusUnauthorized, "Pass not accepted")
		return domain.Actor{}, false
	}

	introduced, inspectErr := usecase.NewInspectPass(g.pool, actor).Exec(r.Context())
	if inspectErr != nil {
		sendProblem(w, gateFailure(inspectErr))
		return domain.Actor{}, false
	}
	if !introduced && !matches(g.nameless, r) {
		writeProblem(w, http.StatusForbidden, "Introduce yourself first")
		return domain.Actor{}, false
	}

	return actor, true
}

// matches says whether r is one of the doors listed, method included: a
// path opened for reading is not thereby opened for writing.
func matches(doors map[string]string, r *http.Request) bool {
	method, listed := doors[r.URL.Path]
	return listed && method == r.Method
}

// gateFailure turns a failure into what the caller is told. It is a type
// switch rather than a chain of checks so that gochecksumtype verifies
// it: a failure added to [domain.InspectPassError] breaks the build here
// instead of quietly becoming a 500.
//
// One wording covers all three ways a pass stops standing. Which one it
// was is not said — that would report to whoever holds a stolen pass
// exactly what became of it — but the remedy is the same in every case,
// and saying nothing at all would leave its rightful holder guessing.
func gateFailure(err domain.InspectPassError) Problem {
	switch err.(type) {
	case domain.PassNotUsableError:
		return NewDetailedProblem(http.StatusUnauthorized, "Sign in again")
	case domain.UnexpectedError:
		return NewProblem(http.StatusInternalServerError)
	}
	return NewProblem(http.StatusInternalServerError)
}

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/andreychh/coopera/internal/domain"
)

// Seal reads back what was stamped. Only the reading half is named here,
// because only the reading half happens at the door; the other half
// belongs to whatever hands passes out, and is named there.
type Seal interface {
	Read(token domain.AccessToken) (domain.Actor, error)
}

// actorKey is unexported and of its own type, so that nothing outside
// this package can write to the place handlers read from. A caller
// cannot forge an actor into the context, because it cannot name the
// key.
type actorKey struct{}

// withActor is the only way an actor gets into a request, and the gate
// is the only caller: by the time a handler is reached, the actor there
// has been through a signature and a look at the session behind it.
func withActor(ctx context.Context, actor domain.Actor) context.Context {
	return context.WithValue(ctx, actorKey{}, actor)
}

// actorFrom takes the actor back out. The second result is false only
// where no gate ran, which is to say on the doors that stand open; a
// handler behind a guarded one can rely on it.
func actorFrom(ctx context.Context) (domain.Actor, bool) {
	actor, present := ctx.Value(actorKey{}).(domain.Actor)
	return actor, present
}

// bearerToken pulls the token out of an Authorization header, and is
// deliberately strict about the shape: exactly the word Bearer, one
// space, and something after it. The word is compared without regard to
// case, as RFC 9110 says a scheme must be.
func bearerToken(r *http.Request) (domain.AccessToken, bool) {
	scheme, token, found := strings.Cut(r.Header.Get("Authorization"), " ")
	if !found || !strings.EqualFold(scheme, "Bearer") || token == "" {
		return "", false
	}
	return domain.AccessToken(token), true
}

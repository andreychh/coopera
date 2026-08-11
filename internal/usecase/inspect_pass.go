// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package usecase

import (
	"context"
	"errors"

	"github.com/andreychh/coopera/internal/db"
	"github.com/andreychh/coopera/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InspectPass answers the two things the door needs about a caller and
// nothing else: whether their pass still stands, and whether its holder
// has a name. It is the one usecase that asks a question instead of
// doing something, and it exists for the door standing in front of
// almost every other one.
//
// A signature says who is asking but cannot say whether the session
// behind them is still open: signing out and catching a theft both close
// one, and neither can reach a token already handed out. That is why
// this question is asked of the database on every guarded request, and
// why asking it makes the answer current rather than as of some minutes
// ago.
//
// Being nameless is not a failure here — it is the answer. The refusal
// it leads to is the door's to give.
type InspectPass struct {
	pool *pgxpool.Pool

	actor domain.Actor
}

func NewInspectPass(pool *pgxpool.Pool, actor domain.Actor) InspectPass {
	return InspectPass{pool: pool, actor: actor}
}

func (u InspectPass) Exec(ctx context.Context) (bool, domain.InspectPassError) {
	username, err := db.New(u.pool).GetPass(ctx, db.GetPassParams{
		UserID:    uuid.UUID(u.actor.ID),
		SessionID: uuid.UUID(u.actor.SessionID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, domain.PassNotUsableError{}
		}
		return false, unexpected("get pass", err)
	}

	// Asked through NewUserState rather than by testing the column, so
	// that what counts as introduced is decided in one place. Should the
	// answer ever depend on more than a name, it changes there and here
	// stays as it is.
	state := domain.NewUserState((*domain.Username)(username))
	_, introduced := state.(domain.Introduced)
	return introduced, nil
}

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

// IsIntroduced answers whether actor has named themselves. It is the one
// usecase that asks a question instead of doing something, and it exists
// for the gate standing in front of every door but two: the gate needs
// the answer to decide, and the meaning of "introduced" belongs here
// rather than with it.
//
// Being nameless is not a failure of this question — it is the answer.
// The refusal it leads to is the gate's to give.
type IsIntroduced struct {
	pool *pgxpool.Pool

	actorID domain.ID
}

func NewIsIntroduced(pool *pgxpool.Pool, actorID domain.ID) IsIntroduced {
	return IsIntroduced{pool: pool, actorID: actorID}
}

func (u IsIntroduced) Exec(ctx context.Context) (bool, domain.IsIntroducedError) {
	row, err := db.New(u.pool).GetUser(ctx, uuid.UUID(u.actorID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, domain.UserNotFoundError{ID: u.actorID}
		}
		return false, unexpected("get user", err)
	}

	// Asked through NewUserState rather than by testing the column, so
	// that what counts as introduced is decided in one place. Should the
	// answer ever depend on more than a name, it changes there and here
	// stays as it is.
	state := domain.NewUserState((*domain.Username)(row.Username))
	_, introduced := state.(domain.Introduced)
	return introduced, nil
}

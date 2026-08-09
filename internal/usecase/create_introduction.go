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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateIntroduction gives actor the name others will know them by. It
// happens once: a person who already has a name is refused rather than
// renamed, because renaming is a different act with a different record
// behind it.
type CreateIntroduction struct {
	pool *pgxpool.Pool

	actorID  domain.ID
	username domain.Username
}

func NewCreateIntroduction(
	pool *pgxpool.Pool,
	actorID domain.ID,
	username domain.Username,
) CreateIntroduction {
	return CreateIntroduction{pool: pool, actorID: actorID, username: username}
}

// Exec needs no transaction: the write is one statement, and it declines
// to touch a person who already has a name, so no two callers can name
// the same one twice.
func (u CreateIntroduction) Exec(
	ctx context.Context,
) (domain.User, domain.CreateIntroductionError) {
	row, err := db.New(u.pool).IntroduceUser(ctx, db.IntroduceUserParams{
		ID:       uuid.UUID(u.actorID),
		Username: new(u.username.String()),
	})
	if err == nil {
		return domain.User{
			ID:        domain.ID(row.ID),
			State:     domain.Introduced{Username: u.username},
			CreatedAt: domain.DateTime(row.CreatedAt),
		}, nil
	}

	pgErr, isPg := errors.AsType[*pgconn.PgError](err)
	if isPg && pgErr.Code == domain.UniqueViolation {
		return domain.User{}, domain.UsernameTakenError{Username: u.username}
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, unexpected("introduce user", err)
	}

	// Nothing was updated, which means either there is no such person or
	// they were introduced long ago. Only the failing path pays for the
	// second question.
	row2, err := db.New(u.pool).GetUser(ctx, uuid.UUID(u.actorID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.UserNotFoundError{ID: u.actorID}
		}
		return domain.User{}, unexpected("get user", err)
	}
	if row2.Username == nil {
		// The person exists and has no name, yet the update declined to
		// give them one — the only thing that declines it. They came into
		// being between the two queries, and neither answer here would be
		// true.
		return domain.User{}, unexpected(
			"introduce user",
			errors.New("user appeared between update and lookup"),
		)
	}
	return domain.User{}, domain.AlreadyIntroducedError{ID: u.actorID}
}

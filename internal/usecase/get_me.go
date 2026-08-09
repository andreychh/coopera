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

// GetMe returns actor as they stand: who they are and whether they have
// introduced themselves. Nothing here can be refused — one's own record
// is not something to be permitted — so the only failure on the caller's
// account is that no such caller exists.
type GetMe struct {
	pool *pgxpool.Pool

	actorID domain.ID
}

func NewGetMe(pool *pgxpool.Pool, actorID domain.ID) GetMe {
	return GetMe{pool: pool, actorID: actorID}
}

func (u GetMe) Exec(ctx context.Context) (domain.User, domain.GetMeError) {
	row, err := db.New(u.pool).GetUser(ctx, uuid.UUID(u.actorID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.UserNotFoundError{ID: u.actorID}
		}
		return domain.User{}, unexpected("get user", err)
	}
	return domain.User{
		ID:        domain.ID(row.ID),
		State:     domain.NewUserState((*domain.Username)(row.Username)),
		CreatedAt: domain.DateTime(row.CreatedAt),
	}, nil
}

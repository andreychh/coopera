// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package usecase

import (
	"context"

	"github.com/andreychh/coopera/internal/db"
	"github.com/andreychh/coopera/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EndSession closes the session the request came from, and only that
// one. The other devices a person is signed in on are untouched: leaving
// a room is not leaving the building.
//
// It is named for what it does rather than for the operation that calls
// it. Nothing is deleted — the row stays and remembers when it ended —
// and the spec says as much in its own words: "ends the session".
type EndSession struct {
	pool *pgxpool.Pool

	actor domain.Actor
}

func NewEndSession(pool *pgxpool.Pool, actor domain.Actor) EndSession {
	return EndSession{pool: pool, actor: actor}
}

// Exec needs no transaction: it is one statement, and that statement
// declines to touch a session already closed, so nothing two callers can
// do together is worse than what one of them does alone.
//
// The door has already asked after this session and found it open, which
// is why nothing here re-checks it. Between that question and this write
// it may close anyway — the same person leaving from two windows at
// once, or a theft caught in between — and the statement is written so
// that this changes nothing worth reporting.
func (u EndSession) Exec(ctx context.Context) domain.EndSessionError {
	err := db.New(u.pool).EndSession(ctx, db.EndSessionParams{
		SessionID: uuid.UUID(u.actor.SessionID),
		UserID:    uuid.UUID(u.actor.ID),
	})
	if err != nil {
		return unexpected("end session", err)
	}
	return nil
}

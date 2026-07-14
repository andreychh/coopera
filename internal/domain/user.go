// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/andreychh/coopera/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User interface {
	CreateTeam(ctx context.Context, name string) (Team, error)
}

type Team struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
}

type SQLUser struct {
	pool *pgxpool.Pool
	id   uuid.UUID
}

func NewSQLUser(pool *pgxpool.Pool, id uuid.UUID) SQLUser {
	return SQLUser{pool: pool, id: id}
}

func (u SQLUser) CreateTeam(ctx context.Context, name string) (Team, error) {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return Team{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	queries := db.New(tx)

	team, err := queries.InsertTeam(ctx, name)
	if err != nil {
		return Team{}, fmt.Errorf("insert team: %w", err)
	}

	_, err = queries.InsertMember(ctx, db.InsertMemberParams{
		TeamID: team.ID,
		UserID: u.id,
		Role:   "owner",
	})
	if err != nil {
		return Team{}, fmt.Errorf("insert owner member: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return Team{}, fmt.Errorf("commit: %w", err)
	}

	return Team{ID: team.ID, Name: team.Name, CreatedAt: team.CreatedAt}, nil
}

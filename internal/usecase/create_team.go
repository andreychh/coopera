// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package usecase

import (
	"context"
	"errors"

	"github.com/andreychh/coopera/internal/db"
	"github.com/andreychh/coopera/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateTeam creates a team and makes actor its owner.
type CreateTeam struct {
	pool *pgxpool.Pool

	actorID domain.ID
	name    domain.TeamName
}

func NewCreateTeam(
	pool *pgxpool.Pool,
	actorID domain.ID,
	name domain.TeamName,
) CreateTeam {
	return CreateTeam{pool: pool, actorID: actorID, name: name}
}

// Exec needs no transaction: the team and its owner are written by one
// statement, so a team without an owner is not a state this can leave
// behind even if the write fails halfway.
func (u CreateTeam) Exec(ctx context.Context) (domain.Team, domain.CreateTeamError) {
	team, err := db.New(u.pool).InsertTeamWithOwner(ctx, db.InsertTeamWithOwnerParams{
		Name:   string(u.name),
		UserID: uuid.UUID(u.actorID),
	})
	if err != nil {
		pgErr, isPg := errors.AsType[*pgconn.PgError](err)
		if isPg && pgErr.Code == domain.ForeignKeyViolation &&
			pgErr.ConstraintName == domain.MembersUserIDForeignKey {
			return domain.Team{}, domain.UserNotFoundError{ID: u.actorID}
		}
		return domain.Team{}, unexpected("insert team with owner", err)
	}

	return domain.Team{
		ID:        domain.ID(team.ID),
		Name:      domain.TeamName(team.Name),
		CreatedAt: domain.DateTime(team.CreatedAt),
	}, nil
}

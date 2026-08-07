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

// ListMyTeams lists the teams actor belongs to, whatever their
// role in each. A team actor has left is not among them: this list and
// GetTeam answer the same question, so a team here is one actor
// can open, and a team missing here is one they cannot.
type ListMyTeams struct {
	pool *pgxpool.Pool

	actorID domain.ID
}

func NewListMyTeams(pool *pgxpool.Pool, actorID domain.ID) ListMyTeams {
	return ListMyTeams{pool: pool, actorID: actorID}
}

// Exec returns an empty slice rather than nil when actor belongs to no
// team: having no teams is an ordinary state, not an absence of answer.
func (u ListMyTeams) Exec(
	ctx context.Context,
) ([]domain.Team, domain.ListMyTeamsError) {
	rows, err := db.New(u.pool).ListTeamsForMember(ctx, uuid.UUID(u.actorID))
	if err != nil {
		return nil, unexpected("list teams for member", err)
	}

	teams := make([]domain.Team, 0, len(rows))
	for _, row := range rows {
		teams = append(teams, domain.Team{
			ID:        domain.ID(row.ID),
			Name:      domain.TeamName(row.Name),
			CreatedAt: domain.DateTime(row.CreatedAt),
		})
	}
	return teams, nil
}

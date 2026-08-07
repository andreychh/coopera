// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package v2

import (
	"context"
	"errors"

	"github.com/andreychh/coopera/internal/db"
	"github.com/andreychh/coopera/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetTeamUsecase returns a team as seen by actor. It fails with
// TeamNotFoundError when the team doesn't exist, when actor never
// belonged to it and when actor has left it, so that none of the three
// can be told from the others: a team is visible only to the people
// currently in it, and past membership grants nothing.
type GetTeamUsecase struct {
	pool *pgxpool.Pool

	actorID domain.ID
	teamID  domain.ID
}

func NewGetTeamUsecase(pool *pgxpool.Pool, actorID, teamID domain.ID) GetTeamUsecase {
	return GetTeamUsecase{pool: pool, actorID: actorID, teamID: teamID}
}

func (u GetTeamUsecase) Exec(ctx context.Context) (domain.Team, domain.GetTeamError) {
	row, err := db.New(u.pool).GetTeamForMember(ctx, db.GetTeamForMemberParams{
		ID:     uuid.UUID(u.teamID),
		UserID: uuid.UUID(u.actorID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Team{}, domain.TeamNotFoundError{ID: u.teamID}
		}
		return domain.Team{}, unexpected("get team for member", err)
	}
	return domain.Team{
		ID:        domain.ID(row.ID),
		Name:      domain.TeamName(row.Name),
		CreatedAt: domain.DateTime(row.CreatedAt),
	}, nil
}

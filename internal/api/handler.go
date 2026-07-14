// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	pool *pgxpool.Pool
}

func NewServer(pool *pgxpool.Pool) Server {
	return Server{pool: pool}
}

func (h Server) CreateTeam(
	ctx context.Context,
	req CreateTeamRequestObject,
) (CreateTeamResponseObject, error) {
	if len(req.Body.Name) < 1 || len(req.Body.Name) > 100 {
		return CreateTeam400ApplicationProblemPlusJSONResponse{
			Title:  http.StatusText(http.StatusBadRequest),
			Status: http.StatusBadRequest,
			Detail: new("name must be between 1 and 100 characters"),
		}, nil
	}

	team, err := domain.NewSQLUser(h.pool, req.Params.XUserId).CreateTeam(ctx, req.Body.Name)
	if err != nil {
		return nil, err
	}

	return CreateTeam201JSONResponse{
		Body: Team{
			Id:        &team.ID,
			Name:      team.Name,
			CreatedAt: &team.CreatedAt,
		},
		Headers: CreateTeam201ResponseHeaders{
			Location: new("/v1/teams/" + team.ID.String()),
		},
	}, nil
}

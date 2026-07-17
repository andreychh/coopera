// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
)

type Server struct {
	world domain.World
}

func NewServer(world domain.World) Server {
	return Server{world: world}
}

func (s Server) CreateTeam(
	ctx context.Context,
	req CreateTeamRequestObject,
) (CreateTeamResponseObject, error) {
	teamName, err := domain.ParseTeamName(req.Body.Name)
	if err != nil {
		//nolint:nilerr // err is translated into a typed 400 response, not propagated.
		return CreateTeam400ApplicationProblemPlusJSONResponse{
			Title:  http.StatusText(http.StatusBadRequest),
			Status: http.StatusBadRequest,
			Detail: new("Invalid team name: " + err.Error()),
		}, nil
	}
	userID, err := domain.ParseID(req.Params.XUserId)
	if err != nil {
		//nolint:nilerr // err is translated into a typed 400 response, not propagated.
		return CreateTeam400ApplicationProblemPlusJSONResponse{
			Title:  http.StatusText(http.StatusBadRequest),
			Status: http.StatusBadRequest,
			Detail: new("Invalid X-User-Id (UUID): " + err.Error()),
		}, nil
	}

	team, err := s.world.User(userID).CreateTeam(ctx, teamName)
	if err != nil {
		if _, ok := errors.AsType[domain.UserNotFoundError](err); ok {
			return CreateTeam401ApplicationProblemPlusJSONResponse{
				Title:  http.StatusText(http.StatusUnauthorized),
				Status: http.StatusUnauthorized,
			}, nil
		}
		return CreateTeam500ApplicationProblemPlusJSONResponse{
			Title:  http.StatusText(http.StatusInternalServerError),
			Status: http.StatusInternalServerError,
		}, nil
	}

	info, err := team.Info(ctx)
	if err != nil {
		//nolint:nilerr // err is translated into a typed 500 response, not propagated.
		return CreateTeam500ApplicationProblemPlusJSONResponse{
			Title:  http.StatusText(http.StatusInternalServerError),
			Status: http.StatusInternalServerError,
		}, nil
	}

	return CreateTeam201JSONResponse{
		Body: Team{
			Id:        new(info.ID.String()),
			Name:      info.Name.String(),
			CreatedAt: new(info.CreatedAt.String()),
		},
		Headers: CreateTeam201ResponseHeaders{
			Location: new("/v1/teams/" + info.ID.String()),
		},
	}, nil
}

func (s Server) GetTeam(
	ctx context.Context,
	req GetTeamRequestObject,
) (GetTeamResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (s Server) ListMyTeams(
	ctx context.Context,
	req ListMyTeamsRequestObject,
) (ListMyTeamsResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

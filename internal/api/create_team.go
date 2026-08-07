// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
)

func (s Server) CreateTeam(
	ctx context.Context,
	req CreateTeamRequestObject,
) (CreateTeamResponseObject, error) {
	teamName, err := domain.ParseTeamName(req.Body.Name)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return CreateTeam400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid team name"),
		), nil
	}
	userID, err := domain.ParseID(req.Params.XUserId)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return CreateTeam400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid X-User-Id"),
		), nil
	}

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	info, createErr := usecase.NewCreateTeam(s.pool, userID, teamName).Exec(ctx)
	if createErr != nil {
		return createTeamError(createErr), nil
	}

	return CreateTeam201JSONResponse{
		Body: newTeam(info),
		Headers: CreateTeam201ResponseHeaders{
			Location: new("/v1/teams/" + info.ID.String()),
		},
	}, nil
}

// createTeamError is a type switch rather than a chain of checks so that
// gochecksumtype verifies it: a failure added to [domain.CreateTeamError]
// breaks the build here instead of quietly becoming a 500.
func createTeamError(err domain.CreateTeamError) CreateTeamResponseObject {
	switch err.(type) {
	case domain.UserNotFoundError:
		return CreateTeam401ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusUnauthorized),
		)
	case domain.UnexpectedError:
		return CreateTeam500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		)
	}
	return CreateTeam500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}

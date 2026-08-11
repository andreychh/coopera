// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
)

func (s Server) GetTeam(
	ctx context.Context,
	req GetTeamRequestObject,
) (GetTeamResponseObject, error) {
	teamID, err := domain.ParseID(req.Id)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return GetTeam400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid id"),
		), nil
	}
	// The gate put the actor here, having read their token and found the
	// session behind it still open. Absence means this handler was reached
	// without the gate in front of it, which is a fault of ours.
	actor, present := actorFrom(ctx)
	if !present {
		return GetTeam500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		), nil
	}

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	info, getErr := usecase.NewGetTeam(s.pool, actor.ID, teamID).Exec(ctx)
	if getErr != nil {
		return getTeamError(getErr), nil
	}

	return GetTeam200JSONResponse(newTeam(info)), nil
}

// getTeamError is a type switch rather than a chain of checks so that
// gochecksumtype verifies it: a failure added to [domain.GetTeamError]
// breaks the build here instead of quietly becoming a 500.
func getTeamError(err domain.GetTeamError) GetTeamResponseObject {
	switch err.(type) {
	case domain.TeamNotFoundError:
		return GetTeam404ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusNotFound),
		)
	case domain.UnexpectedError:
		return GetTeam500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		)
	}
	return GetTeam500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
)

func (s Server) ListMyTeams(
	ctx context.Context,
	req ListMyTeamsRequestObject,
) (ListMyTeamsResponseObject, error) {
	// The gate put the actor here, having read their token and found the
	// session behind it still open. Absence means this handler was reached
	// without the gate in front of it, which is a fault of ours.
	actor, present := actorFrom(ctx)
	if !present {
		return ListMyTeams500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		), nil
	}

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	teams, listErr := usecase.NewListMyTeams(s.pool, actor.ID).Exec(ctx)
	if listErr != nil {
		return listMyTeamsError(listErr), nil
	}

	return ListMyTeams200JSONResponse{Items: newTeams(teams)}, nil
}

// listMyTeamsError is a type switch rather than a chain of checks so
// that gochecksumtype verifies it: a failure added to
// [domain.ListMyTeamsError] breaks the build here instead of quietly
// becoming a 500.
func listMyTeamsError(err domain.ListMyTeamsError) ListMyTeamsResponseObject {
	//nolint:gocritic // an if would not be checked by gochecksumtype, which is the whole point
	switch err.(type) {
	case domain.UnexpectedError:
		return ListMyTeams500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		)
	}
	return ListMyTeams500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
)

func (s Server) ListInviteLinks(
	ctx context.Context,
	req ListInviteLinksRequestObject,
) (ListInviteLinksResponseObject, error) {
	teamID, err := domain.ParseID(req.TeamId)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return ListInviteLinks400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid team_id"),
		), nil
	}
	// The gate put the actor here, having read their token and found the
	// session behind it still open. Absence means this handler was reached
	// without the gate in front of it, which is a fault of ours.
	actor, present := actorFrom(ctx)
	if !present {
		return ListInviteLinks500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		), nil
	}
	var status *domain.LinkStatus
	if req.Params.Status != nil {
		status = new(domain.LinkStatus(*req.Params.Status))
	}

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	links, listErr := usecase.NewListInviteLinks(s.pool, actor.ID, teamID, status).Exec(ctx)
	if listErr != nil {
		return listInviteLinksError(listErr), nil
	}

	body, err := newInviteLinks(links)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return ListInviteLinks500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		), nil
	}

	return ListInviteLinks200JSONResponse{Items: body}, nil
}

// listInviteLinksError is a type switch rather than a chain of checks so
// that gochecksumtype verifies it: a failure added to
// [domain.ListInviteLinksError] breaks the build here instead of quietly
// becoming a 500.
func listInviteLinksError(err domain.ListInviteLinksError) ListInviteLinksResponseObject {
	switch err.(type) {
	case domain.NotTeamOwnerError:
		return ListInviteLinks403ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusForbidden),
		)
	case domain.UnexpectedError:
		return ListInviteLinks500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		)
	}
	return ListInviteLinks500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}

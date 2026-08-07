// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
)

func (s Server) CreateInviteLink(
	ctx context.Context,
	req CreateInviteLinkRequestObject,
) (CreateInviteLinkResponseObject, error) {
	teamID, err := domain.ParseID(req.TeamId)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return CreateInviteLink400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid team_id"),
		), nil
	}
	userID, err := domain.ParseID(req.Params.XUserId)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return CreateInviteLink400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid X-User-Id"),
		), nil
	}

	// An omitted expires_at means the link never expires. That reading of
	// the request belongs here, not in the domain: the domain is told a
	// Validity, it doesn't infer one from a field that wasn't sent.
	var validity domain.Validity = domain.Forever{}
	if req.Body != nil && req.Body.ExpiresAt != nil {
		validity, err = domain.ParseValidity(*req.Body.ExpiresAt)
		if err != nil {
			//nolint:nilerr // outcome is encoded in the response, not the error return
			return CreateInviteLink400ApplicationProblemPlusJSONResponse(
				NewDetailedProblem(http.StatusBadRequest, "Invalid expires_at"),
			), nil
		}
	}

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	link, createErr := usecase.NewCreateInviteLink(s.pool, userID, teamID, validity).Exec(ctx)
	if createErr != nil {
		return createInviteLinkError(createErr), nil
	}

	item, err := newInviteLink(link)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return CreateInviteLink500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		), nil
	}

	return CreateInviteLink201JSONResponse(item), nil
}

// createInviteLinkError is a type switch rather than a chain of
// checks so that gochecksumtype verifies it: a failure added to
// [domain.CreateInviteLinkError] breaks the build here instead of
// quietly becoming a 500.
func createInviteLinkError(err domain.CreateInviteLinkError) CreateInviteLinkResponseObject {
	switch err.(type) {
	case domain.NotTeamOwnerError:
		return CreateInviteLink403ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusForbidden),
		)
	case domain.UnexpectedError:
		return CreateInviteLink500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		)
	}
	return CreateInviteLink500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}

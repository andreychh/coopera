// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
)

func (s Server) RevokeInviteLink(
	ctx context.Context,
	req RevokeInviteLinkRequestObject,
) (RevokeInviteLinkResponseObject, error) {
	userID, err := domain.ParseID(req.Params.XUserId)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return RevokeInviteLink400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid X-User-Id"),
		), nil
	}

	revokeErr := usecase.NewRevokeInviteLink(s.pool, userID, domain.Code(req.Code)).Exec(ctx)
	if revokeErr != nil {
		return revokeInviteLinkError(revokeErr), nil
	}

	return RevokeInviteLink204Response{}, nil
}

// revokeInviteLinkError is a type switch rather than a chain of checks
// so that gochecksumtype verifies it: a failure added to
// [domain.RevokeInviteLinkError] breaks the build here instead of
// quietly becoming a 500.
func revokeInviteLinkError(err domain.RevokeInviteLinkError) RevokeInviteLinkResponseObject {
	switch err.(type) {
	case domain.InviteLinkNotFoundError:
		return RevokeInviteLink404ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusNotFound),
		)
	case domain.NotTeamOwnerError:
		return RevokeInviteLink403ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusForbidden),
		)
	case domain.InviteLinkAlreadyRevokedError:
		return RevokeInviteLink409ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusConflict),
		)
	case domain.UnexpectedError:
		return RevokeInviteLink500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		)
	}
	return RevokeInviteLink500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}

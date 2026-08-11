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
	// The gate put the actor here, having read their token and found the
	// session behind it still open. Absence means this handler was reached
	// without the gate in front of it, which is a fault of ours.
	actor, present := actorFrom(ctx)
	if !present {
		return RevokeInviteLink500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		), nil
	}

	revokeErr := usecase.NewRevokeInviteLink(s.pool, actor.ID, domain.Code(req.Code)).Exec(ctx)
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

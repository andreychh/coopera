// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
)

func (s Server) AcceptInviteLink(
	ctx context.Context,
	req AcceptInviteLinkRequestObject,
) (AcceptInviteLinkResponseObject, error) {
	// The gate put the actor here, having read their token and found the
	// session behind it still open. Absence means this handler was reached
	// without the gate in front of it, which is a fault of ours.
	actor, present := actorFrom(ctx)
	if !present {
		return AcceptInviteLink500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		), nil
	}

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	team, joined, acceptErr := usecase.NewAcceptInviteLink(
		s.pool, actor.ID, domain.Code(req.Code),
	).Exec(ctx)
	if acceptErr != nil {
		return acceptInviteLinkError(acceptErr), nil
	}

	// 201 says a membership came into being; someone who was already in
	// the team gets 200, because nothing did.
	if !joined {
		return AcceptInviteLink200JSONResponse(newTeam(team)), nil
	}
	return AcceptInviteLink201JSONResponse(newTeam(team)), nil
}

// acceptInviteLinkError is a type switch rather than a chain of checks
// so that gochecksumtype verifies it: a failure added to
// [domain.AcceptInviteLinkError] breaks the build here instead of
// quietly becoming a 500.
func acceptInviteLinkError(err domain.AcceptInviteLinkError) AcceptInviteLinkResponseObject {
	switch err.(type) {
	case domain.UserNotFoundError:
		return AcceptInviteLink401ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusUnauthorized),
		)
	case domain.InviteLinkNotFoundError:
		return AcceptInviteLink404ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusNotFound),
		)
	case domain.InviteLinkNotUsableError:
		return AcceptInviteLink410ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusGone),
		)
	case domain.UnexpectedError:
		return AcceptInviteLink500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		)
	}
	return AcceptInviteLink500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}

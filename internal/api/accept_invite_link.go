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
	team, acceptErr := usecase.NewAcceptInviteLink(
		s.pool, actor.ID, domain.Code(req.Code),
	).Exec(ctx)
	if acceptErr != nil {
		return acceptInviteLinkError(acceptErr), nil
	}

	// What came into being is a membership, and a membership has no
	// address of its own, so Location points at the team the new member
	// can now reach.
	return AcceptInviteLink201JSONResponse{
		Body: newTeam(team),
		Headers: AcceptInviteLink201ResponseHeaders{
			Location: new("/v1/teams/" + team.ID.String()),
		},
	}, nil
}

// acceptInviteLinkError is a type switch rather than a chain of checks
// so that gochecksumtype verifies it: a failure added to
// [domain.AcceptInviteLinkError] breaks the build here instead of
// quietly becoming a 500.
func acceptInviteLinkError(err domain.AcceptInviteLinkError) AcceptInviteLinkResponseObject {
	switch e := err.(type) {
	case domain.UserNotFoundError:
		return AcceptInviteLink401ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusUnauthorized),
		)
	case domain.InviteLinkNotFoundError:
		return AcceptInviteLink404ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusNotFound),
		)
	case domain.AlreadyMemberError:
		// The team is named because the caller came holding a code, and a
		// code says nothing about which team stands behind it.
		return AcceptInviteLink409ApplicationProblemPlusJSONResponse{
			TeamId: e.TeamID.String(),
			Detail: new("Already a member of this team"),
			Status: http.StatusConflict,
			Title:  http.StatusText(http.StatusConflict),
		}
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

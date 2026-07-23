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

	team, err := s.world.User(userID).CreateTeam(ctx, teamName)
	if err != nil {
		if _, ok := errors.AsType[domain.UserNotFoundError](err); ok {
			return CreateTeam401ApplicationProblemPlusJSONResponse(
				NewProblem(http.StatusUnauthorized),
			), nil
		}
		return CreateTeam500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		), nil
	}

	info, err := team.Info(ctx)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return CreateTeam500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		), nil
	}

	return CreateTeam201JSONResponse{
		Body: Team{
			Id:        info.ID.String(),
			Name:      info.Name.String(),
			CreatedAt: info.CreatedAt.String(),
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

func (s Server) RevokeInviteLink(
	ctx context.Context,
	req RevokeInviteLinkRequestObject,
) (RevokeInviteLinkResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (s Server) AcceptInviteLink(
	ctx context.Context,
	req AcceptInviteLinkRequestObject,
) (AcceptInviteLinkResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (s Server) ListInviteLinks(
	ctx context.Context,
	req ListInviteLinksRequestObject,
) (ListInviteLinksResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

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

	var expiresAt *domain.InviteLinkExpiry
	if req.Body != nil && req.Body.ExpiresAt != nil {
		var expiry domain.InviteLinkExpiry
		expiry, err = domain.ParseInviteLinkExpiry(*req.Body.ExpiresAt)
		if err != nil {
			//nolint:nilerr // outcome is encoded in the response, not the error return
			return CreateInviteLink400ApplicationProblemPlusJSONResponse(
				NewDetailedProblem(http.StatusBadRequest, "Invalid expires_at"),
			), nil
		}
		expiresAt = &expiry
	}

	_, err = s.world.User(userID).Info(ctx)
	if err != nil {
		return createInviteLinkIdentityError(err), nil
	}

	link, err := s.world.Team(teamID).CreateInviteLink(ctx, userID, expiresAt)
	if err != nil {
		return createInviteLinkActionError(err), nil
	}

	linkState, err := newActiveInviteLinkState(link)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return CreateInviteLink500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		), nil
	}

	return CreateInviteLink201JSONResponse{
		Code:      link.Code.String(),
		CreatedAt: link.CreatedAt.String(),
		State:     linkState,
		UseCount:  int(link.UseCount),
	}, nil
}

func createInviteLinkIdentityError(err error) CreateInviteLinkResponseObject {
	if _, ok := errors.AsType[domain.UserNotFoundError](err); ok {
		return CreateInviteLink401ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusUnauthorized),
		)
	}
	return CreateInviteLink500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}

func createInviteLinkActionError(err error) CreateInviteLinkResponseObject {
	if _, ok := errors.AsType[domain.TeamNotFoundError](err); ok {
		return CreateInviteLink404ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusNotFound),
		)
	}
	if _, ok := errors.AsType[domain.NotTeamOwnerError](err); ok {
		return CreateInviteLink403ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusForbidden),
		)
	}
	return CreateInviteLink500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}

// newActiveInviteLinkState builds the state of an invite link that was
// just created: it can't already be revoked, and ParseInviteLinkExpiry
// already rejected any expiry that isn't in the future, so it's always
// active.
func newActiveInviteLinkState(link domain.InviteLinkInfo) (InviteLinkState, error) {
	active := ActiveInviteLinkState{Status: ActiveInviteLinkStateStatusActive}
	if link.ExpiresAt != nil {
		active.ExpiresAt = new(link.ExpiresAt.String())
	}

	var state InviteLinkState
	err := state.FromActiveInviteLinkState(active)
	return state, err
}

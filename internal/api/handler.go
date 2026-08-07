// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	v2 "github.com/andreychh/coopera/internal/domain/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	pool *pgxpool.Pool
}

func NewServer(pool *pgxpool.Pool) Server {
	return Server{pool: pool}
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

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	info, createErr := v2.NewCreateTeamUsecase(s.pool, userID, teamName).Exec(ctx)
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
	userID, err := domain.ParseID(req.Params.XUserId)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return GetTeam400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid X-User-Id"),
		), nil
	}

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	info, getErr := v2.NewGetTeamUsecase(s.pool, userID, teamID).Exec(ctx)
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

func (s Server) ListMyTeams(
	ctx context.Context,
	req ListMyTeamsRequestObject,
) (ListMyTeamsResponseObject, error) {
	userID, err := domain.ParseID(req.Params.XUserId)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return ListMyTeams400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid X-User-Id"),
		), nil
	}

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	teams, listErr := v2.NewListMyTeamsUsecase(s.pool, userID).Exec(ctx)
	if listErr != nil {
		return listMyTeamsError(listErr), nil
	}

	return ListMyTeams200JSONResponse(newTeams(teams)), nil
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

	revokeErr := v2.NewRevokeInviteLinkUsecase(s.pool, userID, domain.Code(req.Code)).Exec(ctx)
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

func (s Server) AcceptInviteLink(
	ctx context.Context,
	req AcceptInviteLinkRequestObject,
) (AcceptInviteLinkResponseObject, error) {
	userID, err := domain.ParseID(req.Params.XUserId)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return AcceptInviteLink400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid X-User-Id"),
		), nil
	}

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	team, joined, acceptErr := v2.NewAcceptInviteLinkUsecase(
		s.pool, userID, domain.Code(req.Code),
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
	userID, err := domain.ParseID(req.Params.XUserId)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return ListInviteLinks400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid X-User-Id"),
		), nil
	}
	var status *domain.LinkStatus
	if req.Params.Status != nil {
		status = new(domain.LinkStatus(*req.Params.Status))
	}

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	links, listErr := v2.NewListInviteLinksUsecase(s.pool, userID, teamID, status).Exec(ctx)
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

	return ListInviteLinks200JSONResponse(body), nil
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
	link, createErr := v2.NewCreateInviteLinkUsecase(s.pool, userID, teamID, validity).Exec(ctx)
	if createErr != nil {
		return createInviteLinkActionError(createErr), nil
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

// createInviteLinkActionError is a type switch rather than a chain of
// checks so that gochecksumtype verifies it: a failure added to
// [domain.CreateInviteLinkError] breaks the build here instead of
// quietly becoming a 500.
func createInviteLinkActionError(err domain.CreateInviteLinkError) CreateInviteLinkResponseObject {
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

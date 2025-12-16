package http

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpTeams struct {
	userID int64
	client transport.Client
}

func (h httpTeams) All(ctx context.Context) ([]domain.Team, error) {
	resp := findUserResponse{}
	err := h.client.Get(
		ctx,
		transport.URL("users").With("id", strconv.FormatInt(h.userID, 10)).String(),
		&resp,
	)
	if err != nil {
		return nil, fmt.Errorf("getting user %d: %w", h.userID, err)
	}
	teams := make([]domain.Team, 0, len(resp.Teams))
	for _, team := range resp.Teams {
		teams = append(
			teams,
			Team(team.Id, team.Name, h.client),
		)
	}
	return teams, nil
}

func (h httpTeams) Empty(ctx context.Context) (bool, error) {
	resp := findUserResponse{}
	err := h.client.Get(
		ctx,
		transport.URL("users").With("id", strconv.FormatInt(h.userID, 10)).String(),
		&resp,
	)
	if err != nil {
		return false, fmt.Errorf("getting user %d: %w", h.userID, err)
	}
	return len(resp.Teams) == 0, nil
}

func (h httpTeams) TeamWithName(ctx context.Context, name string) (domain.Team, bool, error) {
	resp := findUserResponse{}
	err := h.client.Get(
		ctx,
		transport.URL("users").With("id", strconv.FormatInt(h.userID, 10)).String(),
		&resp,
	)
	if err != nil {
		return nil, false, fmt.Errorf("getting user %d: %w", h.userID, err)
	}
	for _, team := range resp.Teams {
		if team.Name == name {
			return Team(team.Id, team.Name, h.client), true, nil
		}
	}
	return nil, false, nil
}

func Teams(userID int64, client transport.Client) domain.Teams {
	return httpTeams{
		userID: userID,
		client: client,
	}
}

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpTeams struct {
	userID         int64
	userTelegramID int64
	client         transport.Client
}

type usersResponse struct {
	ID    int64 `json:"id"`
	Teams []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"teams"`
}

func (h httpTeams) All(ctx context.Context) ([]domain.Team, error) {
	data, err := h.client.Get(ctx, transport.NewOutcomingURL("users").
		With("telegram_id", strconv.FormatInt(h.userTelegramID, 10)).
		String(),
	)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}
	resp := usersResponse{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	teams := make([]domain.Team, 0, len(resp.Teams))
	for _, t := range resp.Teams {
		teams = append(teams, httpTeam{
			id:     t.ID,
			name:   t.Name,
			client: h.client,
		})
	}
	return teams, nil
}

func (h httpTeams) Empty(ctx context.Context) (bool, error) {
	data, err := h.client.Get(ctx, transport.NewOutcomingURL("users").
		With("telegram_id", strconv.FormatInt(h.userTelegramID, 10)).
		String(),
	)
	if err != nil {
		return false, fmt.Errorf("getting user: %w", err)
	}
	resp := usersResponse{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return false, fmt.Errorf("unmarshaling data: %w", err)
	}
	return len(resp.Teams) == 0, nil
}

func (h httpTeams) TeamWithName(ctx context.Context, name string) (domain.Team, bool, error) {
	data, err := h.client.Get(ctx, transport.NewOutcomingURL("users").
		With("telegram_id", strconv.FormatInt(h.userTelegramID, 10)).
		String(),
	)
	if err != nil {
		return nil, false, fmt.Errorf("getting user: %w", err)
	}
	resp := usersResponse{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, false, fmt.Errorf("unmarshaling data: %w", err)
	}
	for _, t := range resp.Teams {
		if t.Name == name {
			return httpTeam{
				id:     t.ID,
				name:   t.Name,
				client: h.client,
			}, true, nil
		}
	}
	return nil, false, nil
}

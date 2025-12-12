package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpCommunity struct {
	client transport.Client
}

func (h httpCommunity) CreateUser(ctx context.Context, tgID int64, tgUsername string) (domain.User, error) {
	payload, err := json.Marshal(struct {
		ID       int64  `json:"telegram_id"`
		Username string `json:"username"`
	}{tgID, tgUsername})
	if err != nil {
		return nil, fmt.Errorf("marshaling payload: %w", err)
	}
	data, err := h.client.Post(ctx, "users", payload)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	resp := struct {
		ID int64 `json:"id"`
	}{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	return httpUser{
		id:     resp.ID,
		client: h.client,
	}, nil
}

func (h httpCommunity) UserWithTelegramID(ctx context.Context, tgID int64) (domain.User, error) {
	data, err := h.client.Get(
		ctx,
		transport.NewOutcomingURL("users").
			With("telegram_id", strconv.FormatInt(tgID, 10)).
			String(),
	)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}
	resp := struct {
		ID int64 `json:"id"`
	}{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	return httpUser{
		telegramID: tgID,
		id:         resp.ID,
		client:     h.client,
	}, nil
}

func (h httpCommunity) Team(ctx context.Context, id int64) (domain.Team, error) {
	data, err := h.client.Get(
		ctx,
		transport.NewOutcomingURL("teams").
			With("team_id", strconv.FormatInt(id, 10)).
			String(),
	)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", id, err)
	}
	resp := struct {
		Name string `json:"name"`
	}{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	return httpTeam{
		id:     id,
		name:   resp.Name,
		client: h.client,
	}, nil
}

func Community(client transport.Client) domain.Community {
	return httpCommunity{
		client: client,
	}
}

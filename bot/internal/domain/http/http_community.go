package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
	return User(resp.ID, tgID, tgUsername, h.client), nil
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
		ID       int64  `json:"id"`
		Username string `json:"username"`
	}{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	return User(resp.ID, tgID, resp.Username, h.client), nil
}

func (h httpCommunity) UserWithUsername(ctx context.Context, tgUsername string) (domain.User, error) {
	data, err := h.client.Get(
		ctx,
		transport.NewOutcomingURL("users").
			With("username", tgUsername).
			String(),
	)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}
	resp := struct {
		ID         int64 `json:"id"`
		TelegramID int64 `json:"telegram_id"`
	}{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	return User(resp.ID, resp.TelegramID, tgUsername, h.client), nil
}

func (h httpCommunity) UserWithTelegramIDExists(ctx context.Context, tgID int64) (bool, error) {
	_, err := h.client.Get(
		ctx,
		transport.NewOutcomingURL("users").
			With("telegram_id", strconv.FormatInt(tgID, 10)).
			String(),
	)
	var apiErr transport.APIError
	if errors.As(err, &apiErr) {
		if apiErr.StatusCode == http.StatusNotFound {
			return false, nil
		}
	}
	if err != nil {
		return false, fmt.Errorf("getting user: %w", err)
	}
	return true, nil
}

func (h httpCommunity) UserWithUsernameExists(ctx context.Context, tgUsername string) (bool, error) {
	_, err := h.client.Get(
		ctx,
		transport.NewOutcomingURL("users").
			With("username", tgUsername).
			String(),
	)
	var apiErr transport.APIError
	if errors.As(err, &apiErr) {
		if apiErr.StatusCode == http.StatusNotFound {
			return false, nil
		}
	}
	if err != nil {
		return false, fmt.Errorf("getting user: %w", err)
	}
	return true, nil
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
	return Team(id, resp.Name, h.client), nil
}

func Community(client transport.Client) domain.Community {
	return httpCommunity{
		client: client,
	}
}

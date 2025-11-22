package http

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/transport"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
	"github.com/tidwall/gjson"
)

type httpCommunity struct {
	client transport.Client
}

func (h httpCommunity) CreateUser(ctx context.Context, telegramID int64) (domain.User, error) {
	payload, err := json.Object(json.Fields{"telegram_id": json.I64(telegramID)}).Marshal()
	if err != nil {
		return nil, fmt.Errorf("marshaling payload: %w", err)
	}
	data, err := h.client.Post(ctx, "users", payload)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	id := gjson.GetBytes(data, "id")
	if !id.Exists() {
		return nil, fmt.Errorf("field 'id' not found in response")
	}
	return httpUser{
		id:     id.Int(),
		client: h.client,
	}, nil
}

func (h httpCommunity) UserWithTelegramID(telegramID int64) domain.User {
	return &httpUserWithTelegramID{
		telegramID: telegramID,
		client:     h.client,
	}
}

func (h httpCommunity) Team(id int64) domain.Team {
	return httpTeam{
		id:     id,
		client: h.client,
	}
}

func Community(client transport.Client) domain.Community {
	return httpCommunity{
		client: client,
	}
}

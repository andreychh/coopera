package http

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpUser struct {
	id         int64
	tgID       int64
	tgUsername string
	client     transport.Client
}

func (h httpUser) Username() string {
	return h.tgUsername
}

func (h httpUser) ID() int64 {
	return h.id
}

func (h httpUser) CreatedTeams(_ context.Context) (domain.Teams, error) {
	return Teams(h.id, h.tgID, h.client), nil
}

func (h httpUser) CreateTeam(ctx context.Context, name string) (domain.Team, error) {
	payload, err := json.Marshal(struct {
		ID   int64  `json:"user_id"`
		Name string `json:"name"`
	}{h.id, name})
	if err != nil {
		return nil, fmt.Errorf("marshaling payload: %w", err)
	}
	data, err := h.client.Post(ctx, "teams", payload)
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
	return Team(resp.ID, name, h.client), nil
}

func (h httpUser) AssignedTasks(_ context.Context) (domain.Tasks, error) {
	return UserTasks(h.id, h.tgID, h.client), nil
}

func User(id int64, tgID int64, tgUsername string, client transport.Client) domain.User {
	return httpUser{
		id:         id,
		tgID:       tgID,
		tgUsername: tgUsername,
		client:     client,
	}
}

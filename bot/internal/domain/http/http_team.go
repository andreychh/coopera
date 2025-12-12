package http

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpTeam struct {
	id     int64
	name   string
	client transport.Client
}

func (h httpTeam) ID() int64 {
	return h.id
}

func (h httpTeam) Name() string {
	return h.name
}

func (h httpTeam) AddMember(ctx context.Context, user domain.User) (domain.Member, error) {
	payload, err := json.Marshal(struct {
		TeamID int64 `json:"team_id"`
		UserID int64 `json:"user_id"`
	}{h.id, user.ID()})
	if err != nil {
		return nil, fmt.Errorf("marshaling payload: %w", err)
	}
	data, err := h.client.Post(ctx, "memberships", payload)
	if err != nil {
		return nil, fmt.Errorf("adding member: %w", err)
	}
	resp := struct {
		Message string `json:"message"`
	}{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	var id int64
	_, err = fmt.Sscanf(resp.Message, "Member added successfully with id: %d", &id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse id from message: %w", err)
	}
	// TODO: get role from response
	return Member(id, user.Username(), "unknown", h.client), nil
}

func (h httpTeam) Members(_ context.Context) (domain.Members, error) {
	return Members(h.id, h.client), nil
}

func Team(id int64, name string, client transport.Client) domain.Team {
	return httpTeam{
		id:     id,
		name:   name,
		client: client,
	}
}

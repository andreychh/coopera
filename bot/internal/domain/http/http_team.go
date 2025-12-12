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

func (h httpTeam) AddMember(ctx context.Context, userID int64) (domain.Member, error) {
	payload, err := json.Marshal(struct {
		TeamID int64 `json:"team_id"`
		UserID int64 `json:"user_id"`
	}{h.id, userID})
	if err != nil {
		return nil, fmt.Errorf("marshaling payload: %w", err)
	}
	_, err = h.client.Post(ctx, "memberships", payload)
	if err != nil {
		return nil, err
	}
	return httpMember{
		id:         0,
		name:       "",
		dataSource: nil,
	}, nil
}

func (h httpTeam) Members(ctx context.Context) (domain.Members, error) {
	// TODO implement me
	panic("implement me")
}

package http

import (
	"context"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpUser struct {
	id       int64
	username string
	client   transport.Client
}

func (h httpUser) ID() int64 {
	return h.id
}

func (h httpUser) Username() string {
	return h.username
}

func (h httpUser) CreateTeam(ctx context.Context, name string) (domain.Team, error) {
	req := createTeamRequest{
		UserId: h.id,
		Name:   name,
	}
	resp := createTeamResponse{}
	err := h.client.Post(ctx, "teams", req, &resp)
	if err != nil {
		return nil, err
	}
	return Team(resp.Id, resp.Name, h.client), nil
}

func (h httpUser) Teams(ctx context.Context) (domain.Teams, error) {
	return Teams(h.id, h.client), nil
}

func (h httpUser) AssignedTasks(ctx context.Context) (domain.Tasks, error) {
	return UserTasks(h.id, h.username, h.client), nil
}

func User(id int64, username string, client transport.Client) domain.User {
	return httpUser{
		id:       id,
		username: username,
		client:   client,
	}
}

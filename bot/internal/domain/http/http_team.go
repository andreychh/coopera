package http

import (
	"context"

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
	// TODO implement me
	panic("implement me")
}

func (h httpTeam) Members(ctx context.Context) (domain.Members, error) {
	// TODO implement me
	panic("implement me")
}

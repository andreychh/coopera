package http

import (
	"context"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/transport"
)

type httpTeam struct {
	id     int64
	client transport.Client
}

func (h httpTeam) Details(ctx context.Context) (domain.TeamDetails, error) {
	// TODO implement me
	panic("implement me")
}

func (h httpTeam) AddMember(ctx context.Context, user domain.User) (domain.Member, error) {
	// TODO implement me
	panic("implement me")
}

func (h httpTeam) Members(ctx context.Context) ([]domain.Member, error) {
	// TODO implement me
	panic("implement me")
}

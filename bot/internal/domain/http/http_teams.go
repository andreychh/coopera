package http

import (
	"context"
	"iter"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpTeams struct {
	userID int64
	client transport.Client
}

func (h httpTeams) All(ctx context.Context) (iter.Seq2[int64, domain.Team], error) {
	// TODO implement me
	panic("implement me")
}

func (h httpTeams) Empty(ctx context.Context) (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (h httpTeams) TeamWithName(ctx context.Context, name string) (domain.Team, error) {
	// TODO implement me
	panic("implement me")
}

func (h httpTeams) ContainsTeam(ctx context.Context, name string) (bool, error) {
	// TODO implement me
	panic("implement me")
}

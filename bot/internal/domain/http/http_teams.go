package http

import (
	"context"
	"iter"

	"github.com/andreychh/coopera-bot/internal/domain"
)

type httpTeams struct {
}

func (h httpTeams) Details(ctx context.Context) (iter.Seq2[int64, domain.TeamDetails], error) {
	// TODO implement me
	panic("implement me")
}

func (h httpTeams) All(ctx context.Context) (iter.Seq2[int64, domain.Team], error) {
	// TODO implement me
	panic("implement me")
}

func (h httpTeams) Empty(ctx context.Context) (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (h httpTeams) TeamWithName(name string) domain.Team {
	// TODO implement me
	panic("implement me")
}

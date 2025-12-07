package http

import (
	"context"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpMember struct {
	userID     int64
	teamID     int64
	dataSource transport.Client
}

func (h httpMember) Details(ctx context.Context) (domain.MemberDetails, error) {
	// TODO implement me
	panic("implement me")
}

func (h httpMember) CreateTask(ctx context.Context, points int, description string) (domain.Task, error) {
	// TODO implement me
	panic("implement me")
}

func (h httpMember) CreatedTasks(ctx context.Context) ([]domain.Task, error) {
	// TODO implement me
	panic("implement me")
}

package http

import (
	"context"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type cachedMember struct {
	userID int64
	role   string
	client transport.Client
}

func (c cachedMember) Details(ctx context.Context) (domain.MemberDetails, error) {
	panic("not implemented")
}

func (c cachedMember) CreateTask(ctx context.Context, points int, description string) (domain.Task, error) {
	panic("not implemented")
}

func (c cachedMember) CreatedTasks(ctx context.Context) ([]domain.Task, error) {
	panic("not implemented")
}

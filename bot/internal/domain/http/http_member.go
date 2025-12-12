package http

import (
	"context"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpMember struct {
	id         int64
	name       string
	role       string
	dataSource transport.Client
}

func (h httpMember) ID() int64 {
	return h.id
}

func (h httpMember) Name() string {
	return h.name
}

func (h httpMember) Role() string {
	return h.role
}

func (h httpMember) CreateTask(ctx context.Context, points int, title string, description string) (domain.Task, error) {
	// TODO implement me
	panic("implement me")
}

func (h httpMember) CreatedTasks(ctx context.Context) (domain.Tasks, error) {
	// TODO implement me
	panic("implement me")
}

func Member(id int64, name string, role string, dataSource transport.Client) domain.Member {
	return httpMember{
		id:         id,
		name:       name,
		role:       role,
		dataSource: dataSource,
	}
}

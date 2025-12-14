package http

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpMemberTasks struct {
	memberID int64
	client   transport.Client
}

func (h httpMemberTasks) All(ctx context.Context) ([]domain.Task, error) {
	var resp []findTasksResponse
	err := h.client.Get(
		ctx,
		transport.URL("tasks").With("member_id", strconv.FormatInt(h.memberID, 10)).String(),
		&resp,
	)
	if err != nil {
		return nil, fmt.Errorf("getting tasks for member %d: %w", h.memberID, err)
	}
	tasks := make([]domain.Task, 0, len(resp))
	for _, task := range resp {
		tasks = append(tasks, RespToTask(task, h.client))
	}
	return tasks, nil
}

func (h httpMemberTasks) Empty(ctx context.Context) (bool, error) {
	var resp []findTasksResponse
	err := h.client.Get(
		ctx,
		transport.URL("tasks").With("member_id", strconv.FormatInt(h.memberID, 10)).String(),
		&resp,
	)
	if err != nil {
		return false, fmt.Errorf("getting tasks for member %d: %w", h.memberID, err)
	}
	return len(resp) == 0, nil
}

func MemberTasks(memberID int64, client transport.Client) domain.Tasks {
	return httpMemberTasks{
		memberID: memberID,
		client:   client,
	}
}

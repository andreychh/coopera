package http

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpTeamTasks struct {
	teamID int64
	client transport.Client
}

func (h httpTeamTasks) All(ctx context.Context) ([]domain.Task, error) {
	var resp []findTasksResponse
	err := h.client.Get(
		ctx,
		transport.URL("tasks").With("team_id", strconv.FormatInt(h.teamID, 10)).String(),
		&resp,
	)
	if err != nil {
		return nil, fmt.Errorf("getting tasks for team %d: %w", h.teamID, err)
	}
	tasks := make([]domain.Task, 0, len(resp))
	for _, task := range resp {
		tasks = append(tasks, RespToTask(task, h.client))
	}
	return tasks, nil
}

func (h httpTeamTasks) Empty(ctx context.Context) (bool, error) {
	var resp []findTasksResponse
	err := h.client.Get(
		ctx,
		transport.URL("tasks").With("team_id", strconv.FormatInt(h.teamID, 10)).String(),
		&resp,
	)
	if err != nil {
		return false, fmt.Errorf("getting tasks for team %d: %w", h.teamID, err)
	}
	return len(resp) == 0, nil
}

func TeamTasks(teamID int64, client transport.Client) domain.Tasks {
	return httpTeamTasks{
		teamID: teamID,
		client: client,
	}
}

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type task struct {
	ID               int64             `json:"id"`
	TeamID           int64             `json:"team_id"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Points           int               `json:"points"`
	Status           domain.TaskStatus `json:"status"`
	AssignedToMember *int64            `json:"assigned_to_member"`
	CreatedByUser    int64             `json:"created_by_user"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

type httpTeamTasks struct {
	teamID int64
	client transport.Client
}

func (h httpTeamTasks) All(ctx context.Context) ([]domain.Task, error) {
	data, err := h.client.Get(ctx, transport.NewOutcomingURL("tasks").
		With("team_id", strconv.FormatInt(h.teamID, 10)).
		String(),
	)
	if err != nil {
		return nil, fmt.Errorf("getting tasks for team %d: %w", h.teamID, err)
	}
	var resp []task
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	tasks := make([]domain.Task, 0, len(resp))
	for _, t := range resp {
		tasks = append(tasks, Task(t.ID, t.Title, t.Description, t.Points, t.Status, t.TeamID, h.client))
	}
	return tasks, nil
}

func (h httpTeamTasks) Empty(ctx context.Context) (bool, error) {
	data, err := h.client.Get(ctx, transport.NewOutcomingURL("tasks").
		With("team_id", strconv.FormatInt(h.teamID, 10)).
		String(),
	)
	if err != nil {
		return false, fmt.Errorf("getting tasks for team %d: %w", h.teamID, err)
	}
	var resp []task
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return false, fmt.Errorf("unmarshaling data: %w", err)
	}
	return len(resp) == 0, nil
}

func TeamTasks(teamID int64, client transport.Client) domain.Tasks {
	return httpTeamTasks{
		teamID: teamID,
		client: client,
	}
}

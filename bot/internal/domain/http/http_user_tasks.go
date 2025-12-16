package http

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpUserTasks struct {
	userID   int64
	username string
	client   transport.Client
}

func (h httpUserTasks) All(ctx context.Context) ([]domain.Task, error) {
	resp := findUserResponse{}
	err := h.client.Get(
		ctx,
		transport.URL("users").With("id", strconv.FormatInt(h.userID, 10)).String(),
		&resp,
	)
	if err != nil {
		return nil, fmt.Errorf("getting user %s: %w", h.username, err)
	}
	tasks := make([]domain.Task, 0)
	for _, team := range resp.Teams {
		teamTasks, err := h.fetchTasksForTeam(ctx, team.Id)
		if err != nil {
			return nil, fmt.Errorf("fetching tasks for team %d: %w", team.Id, err)
		}
		tasks = append(tasks, teamTasks...)
	}
	return tasks, nil
}

func (h httpUserTasks) fetchTasksForTeam(ctx context.Context, teamID int64) ([]domain.Task, error) {
	memberID, err := h.fetchMemberID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("fetching member ID for team %d: %w", teamID, err)
	}
	var resp []findTasksResponse
	err = h.client.Get(
		ctx,
		transport.URL("tasks").With("team_id", strconv.FormatInt(teamID, 10)).String(),
		&resp,
	)
	if err != nil {
		return nil, fmt.Errorf("getting tasks for team %d: %w", teamID, err)
	}
	tasks := make([]domain.Task, 0, len(resp))
	for _, task := range resp {
		if task.AssignedToMember != nil && *task.AssignedToMember == memberID {
			tasks = append(tasks, RespToTask(task, h.client))
		}
	}
	return tasks, nil
}

func (h httpUserTasks) fetchMemberID(ctx context.Context, teamID int64) (int64, error) {
	var resp findTeamResponse
	err := h.client.Get(
		ctx,
		transport.URL("teams").With("team_id", strconv.FormatInt(teamID, 10)).String(),
		&resp,
	)
	if err != nil {
		return 0, fmt.Errorf("getting team %d: %w", teamID, err)
	}
	for _, member := range resp.Members {
		if member.Username == h.username {
			return member.MemberId, nil
		}
	}
	return 0, fmt.Errorf("member with username %s not found in team %d", h.username, teamID)
}

func (h httpUserTasks) Empty(ctx context.Context) (bool, error) {
	tasks, err := h.All(ctx)
	if err != nil {
		return false, fmt.Errorf("getting all tasks for user %s: %w", h.username, err)
	}
	return len(tasks) == 0, nil
}

func UserTasks(userID int64, username string, client transport.Client) domain.Tasks {
	return httpUserTasks{
		userID:   userID,
		username: username,
		client:   client,
	}
}

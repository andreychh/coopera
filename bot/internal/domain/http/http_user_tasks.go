package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpUserTasks struct {
	userID         int64
	userTelegramID int64
	client         transport.Client
}

func (h httpUserTasks) All(ctx context.Context) ([]domain.Task, error) {
	data, err := h.client.Get(ctx, transport.NewOutcomingURL("users").
		With("telegram_id", strconv.FormatInt(h.userTelegramID, 10)).
		String(),
	)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}
	resp := usersResponse{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling user data: %w", err)
	}
	teamIDs := make([]int64, 0, len(resp.Teams))
	for _, t := range resp.Teams {
		teamIDs = append(teamIDs, t.ID)
	}
	tasks := make([]domain.Task, 0)
	for _, teamID := range teamIDs {
		teamTasks, err := h.fetchTasks(ctx, teamID)
		if err != nil {
			return nil, fmt.Errorf("fetching tasks for team %d: %w", teamID, err)
		}
		tasks = append(tasks, teamTasks...)
	}

	return tasks, nil
}

func (h httpUserTasks) fetchTasks(ctx context.Context, teamID int64) ([]domain.Task, error) {
	myMemberID, err := h.fetchMembershipID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("fetching membership ID for team %d: %w", teamID, err)
	}
	data, err := h.client.Get(ctx, transport.NewOutcomingURL("tasks").
		With("member_id", strconv.FormatInt(myMemberID, 10)).
		String(),
	)
	if err != nil {
		return nil, fmt.Errorf("getting tasks for team %d: %w", teamID, err)
	}
	var resp []task
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling tasks data: %w", err)
	}
	tasks := make([]domain.Task, 0)
	for _, t := range resp {
		tasks = append(tasks, Task(t.ID, t.Title, t.Description, t.Points, t.Status, t.TeamID, h.client))
	}
	return tasks, nil
}

func (h httpUserTasks) fetchMembershipID(ctx context.Context, teamID int64) (int64, error) {
	data, err := h.client.Get(
		ctx,
		transport.NewOutcomingURL("teams").
			With("team_id", strconv.FormatInt(teamID, 10)).
			String(),
	)
	if err != nil {
		return 0, fmt.Errorf("getting team %d: %w", teamID, err)
	}
	resp := struct {
		Members []struct {
			ID     int64 `json:"member_id"`
			UserID int64 `json:"user_id"`
		} `json:"members"`
	}{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return 0, fmt.Errorf("unmarshaling team members: %w", err)
	}
	for _, m := range resp.Members {
		if m.UserID == h.userID {
			return m.ID, nil
		}
	}
	return 0, fmt.Errorf("user %d is not a member of team %d", h.userID, teamID)
}

func (h httpUserTasks) Empty(ctx context.Context) (bool, error) {
	tasks, err := h.All(ctx)
	if err != nil {
		return false, err
	}
	return len(tasks) == 0, nil
}

func UserTasks(userID int64, tgID int64, client transport.Client) domain.Tasks {
	return httpUserTasks{
		userID:         userID,
		userTelegramID: tgID,
		client:         client,
	}
}

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

type getTeamResponse struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedByUser int64     `json:"created_by_user"`
	Members       []struct {
		MemberID int64             `json:"member_id"`
		UserID   int64             `json:"user_id"`
		Username string            `json:"username"`
		Role     domain.MemberRole `json:"role"`
	} `json:"members"`
}

type httpTask struct {
	id          int64
	title       string
	description string
	points      int
	teamID      int64
	status      domain.TaskStatus
	client      transport.Client
}

func (h httpTask) ID() int64 {
	return h.id
}

func (h httpTask) Title() string {
	return h.title
}

func (h httpTask) Points() int {
	return h.points
}

func (h httpTask) Status() domain.TaskStatus {
	return h.status
}

func (h httpTask) Assignee(ctx context.Context) (domain.Member, error) {
	data, err := h.client.Get(
		ctx,
		transport.NewOutcomingURL("teams").
			With("team_id", strconv.FormatInt(h.teamID, 10)).
			String(),
	)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", h.teamID, err)
	}
	var resp getTeamResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	for _, m := range resp.Members {
		if m.MemberID == h.id {
			return Member(m.MemberID, m.UserID, h.teamID, m.Username, m.Role, h.client), nil
		}
	}
	return nil, fmt.Errorf("assignee for task %d not found in team %d", h.id, h.teamID)
}

func (h httpTask) Team(ctx context.Context) (domain.Team, error) {
	data, err := h.client.Get(
		ctx,
		transport.NewOutcomingURL("teams").
			With("team_id", strconv.FormatInt(h.teamID, 10)).
			String(),
	)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", h.teamID, err)
	}
	var resp getTeamResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	return Team(h.teamID, resp.Name, h.client), nil
}

func (h httpTask) Description() string {
	return h.description
}

func (h httpTask) Creator(ctx context.Context) (domain.Member, error) {
	data, err := h.client.Get(
		ctx,
		transport.NewOutcomingURL("tasks").
			With("task_id", strconv.FormatInt(h.id, 10)).
			String(),
	)
	if err != nil {
		return nil, fmt.Errorf("getting task %d: %w", h.id, err)
	}
	var resp []task
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("task %d not found", h.id)
	}
	data, err = h.client.Get(
		ctx,
		transport.NewOutcomingURL("teams").
			With("team_id", strconv.FormatInt(h.teamID, 10)).
			String(),
	)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", h.teamID, err)
	}
	var teamResp getTeamResponse
	err = json.Unmarshal(data, &teamResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	for _, m := range teamResp.Members {
		if m.UserID == resp[0].CreatedByUser {
			return Member(m.MemberID, m.UserID, h.teamID, m.Username, m.Role, h.client), nil
		}
	}
	return nil, fmt.Errorf("creator for task %d not found in team %d", h.id, h.teamID)
}

func (h httpTask) CreatedAt(ctx context.Context) (time.Time, error) {
	data, err := h.client.Get(
		ctx,
		transport.NewOutcomingURL("tasks").
			With("task_id", strconv.FormatInt(h.id, 10)).
			String(),
	)
	if err != nil {
		return time.Time{}, fmt.Errorf("getting task %d: %w", h.id, err)
	}
	var resp []task
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return time.Time{}, fmt.Errorf("unmarshaling data: %w", err)
	}
	if len(resp) == 0 {
		return time.Time{}, fmt.Errorf("task %d not found", h.id)
	}
	return resp[0].CreatedAt, nil
}

func Task(id int64, title string, description string, points int, status domain.TaskStatus, teamID int64, client transport.Client) domain.Task {
	return httpTask{
		id:          id,
		title:       title,
		description: description,
		points:      points,
		teamID:      teamID,
		status:      status,
		client:      client,
	}
}

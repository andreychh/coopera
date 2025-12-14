package http

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpMember struct {
	id     int64
	userID int64
	teamID int64
	name   string
	role   domain.MemberRole
	client transport.Client
}

func (h httpMember) ID() int64 {
	return h.id
}

func (h httpMember) Username() string {
	return h.name
}

func (h httpMember) Role() domain.MemberRole {
	return h.role
}

func (h httpMember) CreateTask(ctx context.Context, title string, description string, points int, assignee domain.Member) (domain.Task, error) {
	req := createTaskRequest{
		TeamId:        h.teamID,
		CurrentUserId: h.userID,
		Title:         title,
		Description:   description,
	}
	if assignee != domain.NullMember() {
		id := assignee.ID()
		req.AssignedToMember = &id
	}
	if points != 0 {
		req.Points = &points
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling payload: %w", err)
	}
	data, err := h.client.Post(ctx, "tasks", payload)
	if err != nil {
		return nil, fmt.Errorf("creating task: %w", err)
	}
	var resp createTaskResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	return Task(resp.ID, resp.Title, resp.Description, resp.Points, resp.Status, resp.TeamID, h.client), nil
}

func (h httpMember) Tasks(ctx context.Context) (domain.Tasks, error) {
	return MemberTasks(h.id, h.client), nil
}

func Member(id int64, userID int64, teamID int64, name string, role domain.MemberRole, client transport.Client) domain.Member {
	return httpMember{
		id:     id,
		userID: userID,
		teamID: teamID,
		name:   name,
		role:   role,
		client: client,
	}
}

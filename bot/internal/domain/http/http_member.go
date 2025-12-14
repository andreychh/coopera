package http

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type createTaskRequest struct {
	TeamId           int64  `json:"team_id"`
	Points           int    `json:"points"`
	CurrentUserId    int64  `json:"current_user_id"`
	AssignedToMember *int64 `json:"assigned_to_member"`
	Title            string `json:"title"`
	Description      string `json:"description"`
}

type createTaskResponse struct {
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
		Points:        points,
		CurrentUserId: h.userID,
		Title:         title,
		Description:   description,
	}
	if assignee != domain.NullMember() {
		id := assignee.ID()
		req.AssignedToMember = &id
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
	return MemberTasks(h.id, h.teamID, h.client), nil
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

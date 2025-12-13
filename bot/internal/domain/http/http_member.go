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
	TeamId        int64  `json:"team_id"`
	Points        int    `json:"points"`
	CurrentUserId int64  `json:"current_user_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
}

type createTaskResponse struct {
	ID               int64     `json:"id"`
	TeamID           int64     `json:"team_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Points           int       `json:"points"`
	Status           string    `json:"status"`
	AssignedToMember int       `json:"assigned_to_member"`
	CreatedByUser    int       `json:"created_by_user"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type httpMember struct {
	id     int64
	userID int64
	teamID int64
	name   string
	role   string
	client transport.Client
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
	payload, err := json.Marshal(createTaskRequest{
		TeamId:        h.teamID,
		Points:        points,
		CurrentUserId: h.userID,
		Title:         title,
		Description:   description,
	})
	if err != nil {
		return nil, fmt.Errorf("marshaling payload: %w", err)
	}
	data, err := h.client.Post(ctx, "memberships", payload)
	if err != nil {
		return nil, fmt.Errorf("adding member: %w", err)
	}
	var resp createTaskResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	return Task(resp.ID, resp.Title, resp.Points, resp.Status, resp.TeamID, h.client), nil
}

func (h httpMember) Tasks(ctx context.Context) (domain.Tasks, error) {
	return MemberTasks(h.id, h.client), nil
}

func Member(id int64, userID int64, teamID int64, name string, role string, client transport.Client) domain.Member {
	return httpMember{
		id:     id,
		userID: userID,
		teamID: teamID,
		name:   name,
		role:   role,
		client: client,
	}
}

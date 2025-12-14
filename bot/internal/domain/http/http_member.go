package http

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpMember struct {
	id       int64
	username string
	role     domain.MemberRole
	teamID   int64
	client   transport.Client
}

func (h httpMember) ID() int64 {
	return h.id
}

func (h httpMember) Username() string {
	return h.username
}

func (h httpMember) Role() domain.MemberRole {
	return h.role
}

func (h httpMember) Stats(ctx context.Context) (domain.MemberStats, error) {
	panic("not implemented")
}

func (h httpMember) CreateDraft(ctx context.Context, title string, description string) (domain.Task, error) {
	userID, err := h.userID(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting user ID for username %s: %w", h.username, err)
	}
	req := createTaskRequest{
		TeamId:           h.teamID,
		Points:           nil,
		CurrentUserId:    userID,
		AssignedToMember: nil,
		Title:            title,
		Description:      description,
	}
	resp := createTaskResponse{}
	err = h.client.Post(ctx, "tasks", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("creating draft task: %w", err)
	}
	return Task(
		resp.Id,
		resp.Title,
		resp.Description,
		nil,
		resp.CreatedAt,
		resp.CreatedByMember,
		nil,
		resp.TeamId,
		domain.StatusDraft,
		h.client,
	), nil
}

func (h httpMember) CreateUnassigned(ctx context.Context, title string, description string, points int) (domain.Task, error) {
	userID, err := h.userID(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting user ID for username %s: %w", h.username, err)
	}
	req := createTaskRequest{
		TeamId:           h.teamID,
		Points:           &points,
		CurrentUserId:    userID,
		AssignedToMember: nil,
		Title:            title,
		Description:      description,
	}
	resp := createTaskResponse{}
	err = h.client.Post(ctx, "tasks", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("creating unassigned task: %w", err)
	}
	return Task(
		resp.Id,
		resp.Title,
		resp.Description,
		resp.Points,
		resp.CreatedAt,
		resp.CreatedByMember,
		nil,
		resp.TeamId,
		domain.StatusOpen,
		h.client,
	), nil
}

func (h httpMember) CreateAssigned(ctx context.Context, title string, description string, points int, memberID int64) (domain.Task, error) {
	userID, err := h.userID(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting user ID for username %s: %w", h.username, err)
	}
	req := createTaskRequest{
		TeamId:           h.teamID,
		Points:           &points,
		CurrentUserId:    userID,
		AssignedToMember: &memberID,
		Title:            title,
		Description:      description,
	}
	resp := createTaskResponse{}
	err = h.client.Post(ctx, "tasks", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("creating unassigned task: %w", err)
	}
	return Task(
		resp.Id,
		resp.Title,
		resp.Description,
		resp.Points,
		resp.CreatedAt,
		resp.CreatedByMember,
		resp.AssignedToMember,
		resp.TeamId,
		domain.StatusInProgress,
		h.client,
	), nil
}

func (h httpMember) AssignedTasks(ctx context.Context) (domain.Tasks, error) {
	return MemberTasks(h.id, h.client), nil
}

func (h httpMember) EstimateTask(ctx context.Context, taskID int64, points int) error {
	userID, err := h.userID(ctx)
	if err != nil {
		return fmt.Errorf("getting user ID for username %s: %w", h.username, err)
	}
	req := updateTaskRequest{
		CurrentUserId:    userID,
		TaskId:           taskID,
		Points:           &points,
		AssignedToMember: nil,
		Title:            nil,
		Description:      nil,
	}
	err = h.client.Patch(ctx, "tasks", req, nil)
	if err != nil {
		return fmt.Errorf("estimating task %d: %w", taskID, err)
	}
	return nil
}

func (h httpMember) AssignTask(ctx context.Context, taskID int64, memberID int64) error {
	userID, err := h.userID(ctx)
	if err != nil {
		return fmt.Errorf("getting user ID for username %s: %w", h.username, err)
	}
	req := updateTaskRequest{
		CurrentUserId:    userID,
		TaskId:           taskID,
		Points:           nil,
		AssignedToMember: &memberID,
		Title:            nil,
		Description:      nil,
	}
	err = h.client.Patch(ctx, "tasks", req, nil)
	if err != nil {
		return fmt.Errorf("assigning task %d to member %d: %w", taskID, memberID, err)
	}
	return nil
}

func (h httpMember) SubmitTaskForReview(ctx context.Context, taskID int64) error {
	userID, err := h.userID(ctx)
	if err != nil {
		return fmt.Errorf("getting user ID for username %s: %w", h.username, err)
	}
	req := updateTaskStatusRequest{
		CurrentUserId: userID,
		TaskId:        taskID,
		Status:        "in_review",
	}
	err = h.client.Patch(ctx, "tasks/status", req, nil)
	if err != nil {
		return fmt.Errorf("submitting task %d for review: %w", taskID, err)
	}
	return nil
}

func (h httpMember) ApproveTask(ctx context.Context, taskID int64) error {
	userID, err := h.userID(ctx)
	if err != nil {
		return fmt.Errorf("getting user ID for username %s: %w", h.username, err)
	}
	req := updateTaskStatusRequest{
		CurrentUserId: userID,
		TaskId:        taskID,
		Status:        "completed",
	}
	err = h.client.Patch(ctx, "tasks/status", req, nil)
	if err != nil {
		return fmt.Errorf("submitting task %d for review: %w", taskID, err)
	}
	return nil
}

func (h httpMember) userID(ctx context.Context) (int64, error) {
	user, exists, err := Community(h.client).UserWithUsername(ctx, h.username)
	if err != nil {
		return 0, fmt.Errorf("getting user with username %s: %w", h.username, err)
	}
	if !exists {
		return 0, fmt.Errorf("user with username %s not found", h.username)
	}
	return user.ID(), nil
}

func Member(
	id int64,
	username string,
	role domain.MemberRole,
	teamID int64,
	client transport.Client,
) domain.Member {
	return httpMember{
		id:       id,
		username: username,
		role:     role,
		teamID:   teamID,
		client:   client,
	}
}

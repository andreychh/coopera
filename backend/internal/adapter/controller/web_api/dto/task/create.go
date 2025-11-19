package task

import (
	"github.com/andreychh/coopera-backend/internal/entity"
)

type CreateTaskRequest struct {
	TeamID        int32  `json:"team_id" validate:"required"`
	Description   string `json:"description" validate:"max=1000"`
	Points        int32  `json:"points" validate:"required,gte=1"`
	CurrentUserID int32  `json:"current_user_id" validate:"required"`
	Title         string `json:"title" validate:"required,min=1,max=255"`
}

type CreateTaskResponse struct {
	ID          int32   `json:"id"`
	TeamID      int32   `json:"team_id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Points      int32   `json:"points"`
	Status      string  `json:"status,omitempty"`
	CreatedBy   int32   `json:"created_by"`
	AssignedTo  *int32  `json:"assigned_to,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
}

func ToEntityCreateTaskRequest(req *CreateTaskRequest) *entity.Task {
	task := &entity.Task{
		TeamID:    req.TeamID,
		Points:    req.Points,
		CreatedBy: req.CurrentUserID,
		Title:     req.Title,
	}

	if req.Description != "" {
		task.Description = &req.Description
	}

	return task
}

func ToCreateTaskResponse(task *entity.Task) *CreateTaskResponse {
	taskResponse := &CreateTaskResponse{
		ID:          task.ID,
		TeamID:      task.TeamID,
		Title:       task.Title,
		Description: task.Description,
		Points:      task.Points,
		CreatedBy:   task.CreatedBy,
		CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if task.Description != nil {
		taskResponse.Description = task.Description
	}

	if task.AssignedTo != nil {
		taskResponse.AssignedTo = task.AssignedTo
	}

	if task.UpdatedAt != nil {
		updatedAt := task.UpdatedAt.Format("2006-01-02T15:04:05Z")
		taskResponse.UpdatedAt = &updatedAt
	}

	return taskResponse
}

package task

import "github.com/andreychh/coopera-backend/internal/entity"

type UpdateTaskRequest struct {
	CurrentUserID int32   `json:"current_user_id" validate:"required"`
	TaskID        int32   `json:"task_id" validate:"required"`
	Title         *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description   *string `json:"description,omitempty"`
	Points        *int32  `json:"points,omitempty"`
	AssignedTo    *int32  `json:"assigned_to,omitempty"`
}

func ToEntityUpdateTaskRequest(req *UpdateTaskRequest) *entity.UpdateTask {
	task := &entity.UpdateTask{
		TaskID: req.TaskID,
	}

	if req.Title != nil {
		task.Title = req.Title
	}

	if req.Description != nil {
		task.Description = req.Description
	}

	if req.Points != nil {
		task.Points = req.Points
	}

	if req.AssignedTo != nil {
		task.AssignedTo = req.AssignedTo
	}

	return task
}

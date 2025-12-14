package task

import "github.com/andreychh/coopera-backend/internal/entity"

type GetTaskRequest struct {
	TaskID   int32 `form:"task_id" validate:"omitempty,gt=0"`
	MemberID int32 `form:"member_id" validate:"omitempty,gt=0"`
	TeamID   int32 `form:"team_id" validate:"omitempty,gt=0"`
}

func ToEntityGetTaskRequest(req *GetTaskRequest) *entity.TaskFilter {
	return &entity.TaskFilter{
		TaskID:   req.TaskID,
		MemberID: req.MemberID,
		TeamID:   req.TeamID,
	}
}

type GetTaskResponse struct {
	ID               int32   `json:"id"`
	TeamID           int32   `json:"team_id"`
	Title            string  `json:"title"`
	Description      *string `json:"description,omitempty"`
	Points           *int32  `json:"points,omitempty"`
	Status           string  `json:"status"`
	AssignedToMember *int32  `json:"assigned_to_member,omitempty"`
	CreatedByUser    int32   `json:"created_by_user"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        *string `json:"updated_at,omitempty"`
}

func ToGetTaskResponse(task *entity.Task) *GetTaskResponse {
	resp := &GetTaskResponse{
		ID:            task.ID,
		TeamID:        task.TeamID,
		Title:         task.Title,
		Status:        task.Status.String(),
		CreatedByUser: task.CreatedBy,
		CreatedAt:     task.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if task.Points != nil {
		resp.Points = task.Points
	}

	if task.Description != nil {
		resp.Description = task.Description
	}

	if task.AssignedToMember != nil {
		resp.AssignedToMember = task.AssignedToMember
	}

	if task.UpdatedAt != nil {
		ts := task.UpdatedAt.Format("2006-01-02T15:04:05Z")
		resp.UpdatedAt = &ts
	}

	return resp
}

func ToGetTaskListResponse(tasks []entity.Task) []*GetTaskResponse {
	result := make([]*GetTaskResponse, 0, len(tasks))
	for i := range tasks {
		result = append(result, ToGetTaskResponse(&tasks[i]))
	}
	return result
}

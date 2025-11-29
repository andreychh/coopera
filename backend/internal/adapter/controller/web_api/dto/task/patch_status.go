package task

import "github.com/andreychh/coopera-backend/internal/entity"

type PatchStatusRequest struct {
	TaskID        int32  `json:"task_id" validate:"required"`
	CurrentUserID int32  `json:"current_user_id" validate:"required"`
	Status        string `json:"status" validate:"required,oneof=open assigned in_review completed archived"`
}

func ToEntityPatchStatusRequest(req *PatchStatusRequest) *entity.TaskStatus {
	return &entity.TaskStatus{
		TaskID:        req.TaskID,
		Status:        req.Status,
		CurrentUserID: req.CurrentUserID,
	}
}

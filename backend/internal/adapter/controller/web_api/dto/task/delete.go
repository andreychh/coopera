package task

type DeleteTaskRequest struct {
	TaskID        int32 `form:"task_id" validate:"required"`
	CurrentUserID int32 `form:"current_user_id" validate:"required"`
}

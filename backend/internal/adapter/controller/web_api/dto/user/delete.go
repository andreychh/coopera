package user

type DeleteUserRequest struct {
	ID int32 `form:"user_id" validate:"required"`
}

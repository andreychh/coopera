package user

import (
	"github.com/andreychh/coopera-backend/internal/entity"
	"time"
)

type CreateUserRequest struct {
	TelegramID int64  `json:"telegram_id" validate:"required"`
	Username   string `json:"username" validate:"required,max=32"`
}

type CreateUserResponse struct {
	ID         int32     `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
}

func FromCreateUserRequest(req *CreateUserRequest) *entity.UserEntity {
	return &entity.UserEntity{
		TelegramID: &req.TelegramID,
		Username:   &req.Username,
	}
}

func ToCreateUserResponse(user *entity.UserEntity) *CreateUserResponse {
	return &CreateUserResponse{
		ID:         *user.ID,
		TelegramID: *user.TelegramID,
		Username:   *user.Username,
		CreatedAt:  user.CreatedAt.Truncate(time.Second),
	}
}

package user

import (
	"github.com/andreychh/coopera-backend/internal/entity"
	"time"
)

type CreateUserRequest struct {
	TelegramID int64 `json:"telegram_id" validate:"required"`
}

type CreateUserResponse struct {
	ID         int32     `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (r *CreateUserRequest) GetTelegramID() int64 { return r.TelegramID }

func ToCreateUserResponse(user *entity.UserEntity) *CreateUserResponse {
	return &CreateUserResponse{
		ID:         *user.ID,
		TelegramID: user.TelegramID,
		CreatedAt:  user.CreatedAt.Truncate(time.Second),
	}
}

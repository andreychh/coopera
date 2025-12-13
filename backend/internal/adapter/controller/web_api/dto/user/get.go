package user

import (
	"fmt"
	"github.com/andreychh/coopera-backend/internal/entity"
)

type GetUserRequest struct {
	ID         int32  `form:"id" validate:"omitempty"`
	TelegramID int64  `form:"telegram_id" validate:"omitempty"`
	UserName   string `form:"username" validate:"omitempty,max=32"`
}

type TeamInfo struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type GetUserResponse struct {
	ID         int32      `json:"id"`
	TelegramID int64      `json:"telegram_id"`
	Username   string     `json:"username"`
	CreatedAt  string     `json:"created_at"`
	Teams      []TeamInfo `json:"teams"`
}

func ToGetUserResponse(user *entity.UserEntity) *GetUserResponse {
	teams := make([]TeamInfo, len(user.Teams))
	for i, t := range user.Teams {
		teams[i] = TeamInfo{
			ID:   t.TeamID,
			Name: t.TeamName,
			Role: string(t.Role),
		}
	}

	return &GetUserResponse{
		ID:         *user.ID,
		TelegramID: *user.TelegramID,
		Username:   *user.Username,
		CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Teams:      teams,
	}
}

func (r *GetUserRequest) Validate() error {
	if r.TelegramID == 0 && r.UserName == "" && r.ID == 0 {
		return fmt.Errorf("either telegram_id or username must or id be provided")
	}
	return nil
}

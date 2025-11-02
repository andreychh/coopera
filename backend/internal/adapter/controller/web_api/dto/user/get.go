package user

import "github.com/andreychh/coopera/internal/entity"

type GetUserRequest struct {
	TelegramID int64 `form:"telegram_id" validate:"required"`
}

type TeamInfo struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type GetUserResponse struct {
	ID         int32      `json:"id"`
	TelegramID int64      `json:"telegram_id"`
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
		TelegramID: user.TelegramID,
		CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Teams:      teams,
	}
}

func (r *GetUserRequest) GetTelegramID() int64 { return r.TelegramID }

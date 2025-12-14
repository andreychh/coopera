package http

import (
	"time"

	"github.com/andreychh/coopera-bot/internal/domain"
)

type createTaskRequest struct {
	TeamId           int64  `json:"team_id"`
	Points           *int   `json:"points"`
	CurrentUserId    int64  `json:"current_user_id"`
	AssignedToMember *int64 `json:"assigned_to_member"`
	Title            string `json:"title"`
	Description      string `json:"description"`
}

type createTaskResponse struct {
	ID               int64             `json:"id"`
	TeamID           int64             `json:"team_id"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Points           int               `json:"points"`
	Status           domain.TaskStatus `json:"status"`
	AssignedToMember *int64            `json:"assigned_to_member"`
	CreatedByUser    int64             `json:"created_by_user"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

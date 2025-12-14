package http

import (
	"fmt"
	"time"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type createTaskRequest struct {
	CurrentUserId    int64  `json:"current_user_id"`
	TeamId           int64  `json:"team_id"`
	Points           *int   `json:"points"`
	AssignedToMember *int64 `json:"assigned_to_member"`
	Title            string `json:"title"`
	Description      string `json:"description"`
}

type createTaskResponse struct {
	Id               int64     `json:"id"`
	TeamId           int64     `json:"team_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Points           *int      `json:"points"`
	Status           string    `json:"status"`
	CreatedByUser    int64     `json:"created_by_user"`
	CreatedByMember  int64     `json:"created_by_member"`
	AssignedToMember *int64    `json:"assigned_to_member"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type createUserRequest struct {
	TelegramId int64  `json:"telegram_id"`
	Username   string `json:"username"`
}

type createUserResponse struct {
	Id         int64     `json:"id"`
	TelegramId int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
}

type findUserResponse struct {
	Id         int64     `json:"id"`
	TelegramId int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
	Teams      []struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
		Role string `json:"role"`
	} `json:"teams"`
}

type findTeamResponse struct {
	Id            int64     `json:"id"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedByUser int64     `json:"created_by_user"`
	Members       []struct {
		MemberId int64  `json:"member_id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	} `json:"members"`
}

type createMemberRequest struct {
	TeamId int64 `json:"team_id"`
	UserId int64 `json:"user_id"`
}

type createMemberResponse struct {
	Id int64 `json:"id"`
}

type findTasksResponse struct {
	Id               int64     `json:"id"`
	TeamId           int64     `json:"team_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Points           *int      `json:"points"`
	Status           string    `json:"status"`
	AssignedToMember *int64    `json:"assigned_to_member"`
	CreatedByUser    int64     `json:"created_by_user"`
	CreatedByMember  int64     `json:"created_by_member"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type updateTaskRequest struct {
	CurrentUserId    int64   `json:"current_user_id"`
	TaskId           int64   `json:"task_id"`
	Points           *int    `json:"points"`
	AssignedToMember *int64  `json:"assigned_to_member"`
	Title            *string `json:"title"`
	Description      *string `json:"description"`
}

type updateTaskStatusRequest struct {
	CurrentUserId int64  `json:"current_user_id"`
	TaskId        int64  `json:"task_id"`
	Status        string `json:"status"`
}

type createTeamRequest struct {
	UserId int64  `json:"user_id"`
	Name   string `json:"name"`
}

type createTeamResponse struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int64     `json:"created_by"`
}

func RespToTask(resp findTasksResponse, client transport.Client) domain.Task {
	var status domain.TaskStatus
	switch {
	case resp.Points == nil:
		status = domain.StatusDraft
	case resp.AssignedToMember == nil:
		status = domain.StatusOpen
	case resp.Status == "assigned":
		status = domain.StatusInProgress
	case resp.Status == "in_review":
		status = domain.StatusInReview
	case resp.Status == "completed":
		status = domain.StatusDone
	default:
		panic(fmt.Sprintf("unknown task status: %q", resp))
	}
	return Task(
		resp.Id,
		resp.Title,
		resp.Description,
		resp.Points,
		resp.CreatedAt,
		resp.CreatedByMember,
		resp.AssignedToMember,
		resp.TeamId,
		status,
		client,
	)
}

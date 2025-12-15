package domain

import (
	"context"
	"time"
)

type TaskStatus string

const (
	StatusDraft      TaskStatus = "draft"
	StatusOpen       TaskStatus = "open"
	StatusInProgress TaskStatus = "in_progress"
	StatusInReview   TaskStatus = "in_review"
	StatusDone       TaskStatus = "completed"
)

type MemberRole string

const (
	RoleManager MemberRole = "manager"
	RoleMember  MemberRole = "member"
)

type Community interface {
	CreateUser(ctx context.Context, tgID int64, tgUsername string) (User, error)

	UserWithID(ctx context.Context, id int64) (User, bool, error)
	UserWithTelegramID(ctx context.Context, tgID int64) (User, bool, error)
	UserWithUsername(ctx context.Context, tgUsername string) (User, bool, error)
	Team(ctx context.Context, id int64) (Team, bool, error)
	Task(ctx context.Context, id int64) (Task, bool, error)
}

type User interface {
	ID() int64
	Username() string

	CreateTeam(ctx context.Context, name string) (Team, error)

	Teams(ctx context.Context) (Teams, error)
	AssignedTasks(ctx context.Context) (Tasks, error)
}

type Teams interface {
	All(ctx context.Context) ([]Team, error)
	Empty(ctx context.Context) (bool, error)
	TeamWithName(ctx context.Context, name string) (Team, bool, error)
}

type Team interface {
	ID() int64
	Name() string

	Stats(ctx context.Context) (TeamStats, error)

	AddMember(ctx context.Context, userID int64) (Member, error)

	Members(ctx context.Context) (Members, error)
	Tasks(ctx context.Context) (Tasks, error)
}

type TeamStats struct {
	TotalTasks     int
	CompletedTasks int
	TotalPoints    int
}

type Members interface {
	All(ctx context.Context) ([]Member, error)
	MemberWithUsername(ctx context.Context, username string) (Member, bool, error)
}

type Member interface {
	ID() int64
	Username() string
	Role() MemberRole

	Stats(ctx context.Context) (MemberStats, error)

	CreateDraft(ctx context.Context, title string, description string) (Task, error)
	CreateUnassigned(ctx context.Context, title string, description string, points int) (Task, error)
	CreateAssigned(ctx context.Context, title string, description string, points int, memberID int64) (Task, error)
	AssignedTasks(ctx context.Context) (Tasks, error)

	EstimateTask(ctx context.Context, taskID int64, points int) error
	AssignTask(ctx context.Context, taskID int64, memberID int64) error
	SubmitTaskForReview(ctx context.Context, taskID int64) error
	ApproveTask(ctx context.Context, taskID int64) error
}

type MemberStats struct {
	CompletedTasks int
	TotalPoints    int
}

type Tasks interface {
	All(ctx context.Context) ([]Task, error)
	Empty(ctx context.Context) (bool, error)
}

type Task interface {
	ID() int64
	Title() string
	Description() string
	Points() (int, bool)
	Status() TaskStatus
	CreatedAt() time.Time

	CreatedBy(ctx context.Context) (Member, error)
	Assignee(ctx context.Context) (Member, bool, error)
	Team(ctx context.Context) (Team, error)
}

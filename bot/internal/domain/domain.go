package domain

import (
	"context"
)

type TaskStatus string

const (
	StatusOpen      = "open"
	StatusAssigned  = "assigned"
	StatusCompleted = "completed"
)

type MemberRole string

const (
	RoleManager = "manager"
	RoleMember  = "member"
)

type Community interface {
	CreateUser(ctx context.Context, tgID int64, tgUsername string) (User, error)
	UserWithTelegramID(ctx context.Context, tgID int64) (User, error)
	UserWithUsername(ctx context.Context, tgUsername string) (User, error)
	Team(ctx context.Context, id int64) (Team, error)
	UserWithTelegramIDExists(ctx context.Context, tgID int64) (bool, error)
	UserWithUsernameExists(ctx context.Context, tgUsername string) (bool, error)
	Task(ctx context.Context, id int64) (Task, error)
}

type User interface {
	ID() int64
	Username() string
	CreatedTeams(ctx context.Context) (Teams, error)
	CreateTeam(ctx context.Context, name string) (Team, error)
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
	AddMember(ctx context.Context, user User) (Member, error)
	Members(ctx context.Context) (Members, error)
	MemberWithUserID(ctx context.Context, id int64) (Member, error)
	ContainsUser(ctx context.Context, user User) (bool, error)
	Tasks(ctx context.Context) (Tasks, error)
}

type Members interface {
	All(ctx context.Context) ([]Member, error)
}

type Member interface {
	ID() int64
	Name() string
	Role() MemberRole
	CreateTask(ctx context.Context, title string, description string, points int, assignee Member) (Task, error)
	Tasks(ctx context.Context) (Tasks, error)
}

type Tasks interface {
	All(ctx context.Context) ([]Task, error)
	Empty(ctx context.Context) (bool, error)
}

type Task interface {
	ID() int64
	Title() string
	Points() int
	Status() TaskStatus
	Assignee(ctx context.Context) (Member, error)
	Team(ctx context.Context) (Team, error)
}

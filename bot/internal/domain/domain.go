package domain

import (
	"context"
)

type Community interface {
	CreateUser(ctx context.Context, telegramID int64) (User, error)
	UserWithTelegramID(telegramID int64) User
	Team(id int64) Team
}

type User interface {
	CreatedTeams(ctx context.Context) ([]Team, error)
	CreateTeam(ctx context.Context, name string) (Team, error)
}

type Team interface {
	Details(ctx context.Context) (TeamDetails, error)
	AddMember(ctx context.Context, user User) (Member, error)
	Members(ctx context.Context) ([]Member, error)
}

type TeamDetails interface {
	ID() int64
	Name() string
}

type Member interface {
	Details(ctx context.Context) (MemberDetails, error)
	CreateTask(ctx context.Context, points int, description string) (Task, error)
	CreatedTasks(ctx context.Context) ([]Task, error)
}

type MemberDetails interface {
	ID() int64
	Name() string
}

type Task interface {
}

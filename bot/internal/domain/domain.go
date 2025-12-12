package domain

import (
	"context"
)

type Community interface {
	CreateUser(ctx context.Context, tgID int64, tgUsername string) (User, error)
	UserWithTelegramID(ctx context.Context, tgID int64) (User, error)
	Team(ctx context.Context, id int64) (Team, error)
}

type User interface {
	ID() int64
	CreatedTeams(ctx context.Context) (Teams, error)
	CreateTeam(ctx context.Context, name string) (Team, error)
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
}

type Members interface {
	All(ctx context.Context) ([]Member, error)
	Empty(ctx context.Context) (bool, error)
}

type Member interface {
	ID() int64
	Name() string
	CreateTask(ctx context.Context, points int, title string, description string) (Task, error)
	CreatedTasks(ctx context.Context) (Tasks, error)
}

type Tasks interface{}

type Task interface{}

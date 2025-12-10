package domain

import (
	"context"
	"iter"
)

type Community interface {
	CreateUser(ctx context.Context, tgID int64, tgUsername string) (User, error)
	UserWithTelegramID(tgID int64) User
	Team(id int64) Team
}

type User interface {
	Details(ctx context.Context) (UserDetails, error)
	CreatedTeams() Teams
	CreateTeam(ctx context.Context, name string) (Team, error)
}

type UserDetails interface {
	ID() int64
}

type Teams interface {
	Details(ctx context.Context) (iter.Seq2[int64, TeamDetails], error)
	All(ctx context.Context) (iter.Seq2[int64, Team], error)
	Empty(ctx context.Context) (bool, error)
	TeamWithName(name string) Team
}

type Team interface {
	Exists(ctx context.Context) (bool, error)
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

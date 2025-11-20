// This file implements an in-memory version of the Community interface for testing or prototyping purposes.
package memory

import (
	"context"

	"github.com/andreychh/coopera-bot/internal/domain"
)

type MemoryCommunity struct {
}

func (m MemoryCommunity) Team(id int64) domain.Team {
	return memoryTeam{id: id, name: "teamName"}
}

func (m MemoryCommunity) CreateUser(ctx context.Context, telegramID int64) (domain.User, error) {
	return memoryUser{}, nil
}

func (m MemoryCommunity) UserWithTelegramID(telegramID int64) domain.User {
	return memoryUser{}
}

type memoryUser struct {
}

func (m memoryUser) CreatedTeams(ctx context.Context) ([]domain.Team, error) {
	return []domain.Team{
		memoryTeam{1, "teamOne"},
		memoryTeam{2, "teamTwo"},
	}, nil
}

func (m memoryUser) CreateTeam(ctx context.Context, name string) (domain.Team, error) {
	return memoryTeam{123, "createdTeam"}, nil
}

type memoryTeam struct {
	id   int64
	name string
}

func (m memoryTeam) Details(ctx context.Context) (domain.TeamDetails, error) {
	return memoryTeamDetails(m), nil
}

func (m memoryTeam) AddMember(ctx context.Context, user domain.User) (domain.Member, error) {
	// TODO implement me
	panic("implement me")
}

func (m memoryTeam) Members(ctx context.Context) ([]domain.Member, error) {
	// TODO implement me
	panic("implement me")
}

type memoryTeamDetails struct {
	id   int64
	name string
}

func (m memoryTeamDetails) ID() int64 {
	return m.id
}

func (m memoryTeamDetails) Name() string {
	return m.name
}

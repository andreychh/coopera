package usecase

import (
	"context"
	"github.com/andreychh/coopera-backend/internal/entity"
)

type UserUseCase interface {
	CreateUsecase(ctx context.Context, euser entity.UserEntity) (entity.UserEntity, error)
	GetUsecase(ctx context.Context, opts ...any) (entity.UserEntity, error)
	DeleteUsecase(ctx context.Context, userID int32) error
}

type TeamUseCase interface {
	CreateUsecase(ctx context.Context, team entity.TeamEntity) (entity.TeamEntity, error)
	DeleteUsecase(ctx context.Context, teamID, currentUserID int32) error
	GetByIDUsecase(ctx context.Context, teamID int32) (entity.TeamEntity, []entity.MembershipEntity, map[int32]string, error)
	ExistTeamByIDUsecase(ctx context.Context, teamID int32) (bool, error)
}

type MembershipUseCase interface {
	AddMemberUsecase(ctx context.Context, membership entity.MembershipEntity) (int32, error)
	DeleteMemberUsecase(ctx context.Context, memberID, teamID, currentUserID int32) error
	GetMembersUsecase(ctx context.Context, teamID int32) ([]entity.MembershipEntity, error)
	ExistsMemberUsecase(ctx context.Context, memberID int32) (bool, error)
	GetMemberUsecase(ctx context.Context, teamID, memberID int32) (entity.MembershipEntity, error)
}

type TaskUseCase interface {
	CreateUsecase(ctx context.Context, task entity.Task) (entity.Task, error)
	GetUsecase(ctx context.Context, taskFilter entity.TaskFilter) ([]entity.Task, error)
	UpdateStatus(ctx context.Context, status entity.TaskStatus) error
	DeleteUsecase(ctx context.Context, taskID, currentUserID int32) error
	UpdateUsecase(ctx context.Context, task entity.UpdateTask, currentUserID int32) error
	UpdateStatusForEngine(ctx context.Context, taskStatus entity.TaskStatus) error
	GetAllTasks(ctx context.Context) ([]entity.Task, error)
	UpdateForEngine(ctx context.Context, task entity.UpdateTask) error
}

type TaskAssignmentUsecase interface {
	AssignTasks(ctx context.Context) error
}

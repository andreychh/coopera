package usecase

import (
	"context"
	"github.com/andreychh/coopera/internal/entity"
)

// Здесь будут заданы интерфейсы для имплементации конкретных юзкейсов
type UserUseCase interface {
	CreateUsecase(ctx context.Context, euser entity.UserEntity) (entity.UserEntity, error)
	GetUsecase(ctx context.Context, euser entity.UserEntity) (entity.UserEntity, error)
}

type TeamUseCase interface {
	CreateUsecase(ctx context.Context, team entity.TeamEntity) (entity.TeamEntity, error)
	GetByIDUsecase(ctx context.Context, teamID int32) (entity.TeamEntity, []entity.MembershipEntity, error)
}

type MembershipUseCase interface {
	AddMemberUsecase(ctx context.Context, membership entity.MembershipEntity) error
	GetMembersUsecase(ctx context.Context, teamID int32) ([]entity.MembershipEntity, error)
}

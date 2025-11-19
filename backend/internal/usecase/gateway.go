package usecase

import (
	"context"
	"github.com/andreychh/coopera-backend/internal/entity"
	"github.com/jackc/pgx/v4"
)

// Здесь будут заданы интерфейсы для взаимодействия с имплементациями сдоя репозитория и управления транзакциями

type TransactionManageRepository interface {
	WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
	WithinTransactionWithIsolation(ctx context.Context, level pgx.TxIsoLevel, fn func(txCtx context.Context) error) error
}

type UserRepository interface {
	CreateRepo(ctx context.Context, euser entity.UserEntity) (entity.UserEntity, error)
	DeleteRepo(ctx context.Context, userID int32) error
	GetByTelegramIDRepo(ctx context.Context, telegramID int64) (entity.UserEntity, error)
}

type TeamRepository interface {
	CreateRepo(ctx context.Context, team entity.TeamEntity) (entity.TeamEntity, error)
	DeleteRepo(ctx context.Context, teamID int32) error
	GetByIDRepo(ctx context.Context, teamID int32) (entity.TeamEntity, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
}

type MembershipRepository interface {
	AddMemberRepo(ctx context.Context, membership entity.MembershipEntity) error
	DeleteMemberRepo(ctx context.Context, membership entity.MembershipEntity) error
	GetMembersRepo(ctx context.Context, teamID int32) ([]entity.MembershipEntity, error)
}

type TaskRepository interface {
	CreateRepo(ctx context.Context, task entity.Task) (entity.Task, error)
}

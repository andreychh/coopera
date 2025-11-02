package memberships

import (
	"context"
	"github.com/andreychh/coopera/internal/entity"
	"github.com/andreychh/coopera/internal/usecase"
)

type MembershipsUsecase struct {
	txManager            usecase.TransactionManager
	membershipRepository usecase.MembershipRepository
}

func NewMembershipsUsecase(mRepo usecase.MembershipRepository, txManager usecase.TransactionManager) *MembershipsUsecase {
	return &MembershipsUsecase{
		txManager:            txManager,
		membershipRepository: mRepo,
	}
}

func (uc *MembershipsUsecase) AddMemberUsecase(ctx context.Context, membership entity.MembershipEntity) error {
	if membership.Role == "" {
		membership.Role = entity.RoleMember
	}

	// Оборачиваем вызов репозитория в транзакцию
	return uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		return uc.membershipRepository.AddMemberRepo(txCtx, membership)
	})
}

func (uc *MembershipsUsecase) GetMembersUsecase(ctx context.Context, teamID int32) ([]entity.MembershipEntity, error) {
	return uc.membershipRepository.GetMembersRepo(ctx, teamID)
}

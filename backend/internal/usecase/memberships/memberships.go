package memberships

import (
	"context"
	"fmt"
	"github.com/andreychh/coopera-backend/internal/entity"
	"github.com/andreychh/coopera-backend/internal/usecase"
	appErr "github.com/andreychh/coopera-backend/pkg/errors"
)

type MembershipsUsecase struct {
	txManager            usecase.TransactionManageRepository
	membershipRepository usecase.MembershipRepository
}

func NewMembershipsUsecase(memberRepo usecase.MembershipRepository, txManager usecase.TransactionManageRepository) *MembershipsUsecase {
	return &MembershipsUsecase{
		txManager:            txManager,
		membershipRepository: memberRepo,
	}
}

func (uc *MembershipsUsecase) AddMemberUsecase(ctx context.Context, membership entity.MembershipEntity) (int32, error) {
	if membership.Role == "" {
		membership.Role = entity.RoleMember
	}

	var memberID int32
	var repoErr error

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		memberID, repoErr = uc.membershipRepository.AddMemberRepo(txCtx, membership)
		if repoErr != nil {
			return fmt.Errorf("failed to add member: %w", repoErr)
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return memberID, nil
}

func (uc *MembershipsUsecase) DeleteMemberUsecase(ctx context.Context, membership entity.MembershipEntity, currentUserID int32) error {
	return uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		members, err := uc.membershipRepository.GetMembersRepo(txCtx, membership.TeamID)
		if err != nil {
			return fmt.Errorf("failed to get members: %w", err)
		}

		var (
			managerID int32
			found     bool
		)

		for _, m := range members {
			if m.Role == entity.RoleManager {
				managerID = m.UserID
			}
			if m.UserID == membership.UserID {
				found = true
			}
		}

		if !found {
			return appErr.ErrNotFound
		}

		if currentUserID != managerID && currentUserID != membership.UserID {
			return appErr.ErrNoPermissionToDelete
		}

		if membership.UserID == managerID && currentUserID == managerID {
			return appErr.ErrUserOwner
		}

		if err := uc.membershipRepository.DeleteMemberRepo(txCtx, membership); err != nil {
			return fmt.Errorf("failed to delete member: %w", err)
		}

		return nil
	})
}

func (uc *MembershipsUsecase) GetMembersUsecase(ctx context.Context, teamID int32) ([]entity.MembershipEntity, error) {
	var members []entity.MembershipEntity

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		var err error
		members, err = uc.membershipRepository.GetMembersRepo(txCtx, teamID)
		if err != nil {
			return fmt.Errorf("failed to get members: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (uc *MembershipsUsecase) ExistsMemberUsecase(ctx context.Context, id int32) (bool, error) {
	var exists bool
	var err error

	err = uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		exists, err = uc.membershipRepository.MemberExistsRepo(txCtx, id)
		if err != nil {
			return fmt.Errorf("failed to get members: %w", err)
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (uc *MembershipsUsecase) GetMemberUsecase(ctx context.Context, teamID, memberID int32) (entity.MembershipEntity, error) {
	var membership entity.MembershipEntity
	var err error

	err = uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		membership, err = uc.membershipRepository.GetMemberRepo(txCtx, teamID, memberID)
		if err != nil {
			return fmt.Errorf("failed to get member: %w", err)
		}

		return nil
	})
	if err != nil {
		return entity.MembershipEntity{}, err
	}

	return membership, nil
}

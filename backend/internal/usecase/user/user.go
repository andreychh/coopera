package user

import (
	"context"
	"fmt"
	"github.com/andreychh/coopera-backend/internal/entity"
	"github.com/andreychh/coopera-backend/internal/usecase"
)

type UserUsecase struct {
	txManager      usecase.TransactionManageRepository
	userRepository usecase.UserRepository
}

func NewUserUsecase(userRepository usecase.UserRepository, txManager usecase.TransactionManageRepository) *UserUsecase {
	return &UserUsecase{
		txManager:      txManager,
		userRepository: userRepository,
	}
}

func (uc *UserUsecase) CreateUsecase(ctx context.Context, euser entity.UserEntity) (entity.UserEntity, error) {
	var createdUser entity.UserEntity

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		var err error
		createdUser, err = uc.userRepository.CreateRepo(txCtx, euser)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		return nil
	})
	if err != nil {
		return entity.UserEntity{}, err
	}

	return createdUser, nil
}

func (uc *UserUsecase) GetUsecase(ctx context.Context, opts ...any) (entity.UserEntity, error) {
	var user entity.UserEntity

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		var err error
		user, err = uc.userRepository.GetRepo(txCtx, opts...)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		return nil
	})
	if err != nil {
		return entity.UserEntity{}, err
	}

	return user, nil
}

func (uc *UserUsecase) DeleteUsecase(ctx context.Context, userID int32) error {
	return uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.userRepository.DeleteRepo(txCtx, userID); err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
		return nil
	})
}

package user_repo

import (
	"context"
	"github.com/andreychh/coopera-backend/internal/adapter/repository/converter"
	"github.com/andreychh/coopera-backend/internal/adapter/repository/postgres/dao"
	"github.com/andreychh/coopera-backend/internal/entity"
)

type UserRepository struct {
	UserDAO dao.UserRepository
}

func NewUserRepository(userDAO dao.UserRepository) *UserRepository {
	return &UserRepository{
		UserDAO: userDAO,
	}
}

func (ur *UserRepository) CreateRepo(ctx context.Context, euser entity.UserEntity) (entity.UserEntity, error) {
	userModel := converter.FromEntityToModel(euser)
	enuser, err := ur.UserDAO.Create(ctx, userModel)
	if err != nil {
		return entity.UserEntity{}, err
	}

	return enuser, nil
}

func (ur *UserRepository) DeleteRepo(ctx context.Context, userID int32) error {
	return ur.UserDAO.Delete(ctx, userID)
}

func (ur *UserRepository) GetByTelegramIDRepo(ctx context.Context, telegramID int64) (entity.UserEntity, error) {
	return ur.UserDAO.GetByTelegramID(ctx, telegramID)
}

package user_repo

import (
	"context"
	"github.com/andreychh/coopera-backend/internal/adapter/repository/converter"
	repoErr "github.com/andreychh/coopera-backend/internal/adapter/repository/errors"
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

func (ur *UserRepository) GetRepo(ctx context.Context, opts ...any) (entity.UserEntity, error) {
	var (
		telegramID int64
		username   string
		user       entity.UserEntity
		err        error
	)

	for _, opt := range opts {
		switch v := opt.(type) {
		case int64:
			telegramID = v
		case string:
			username = v
		}
	}

	switch {
	case telegramID != 0:
		user, err = ur.UserDAO.GetByTelegramID(ctx, telegramID)
	case username != "":
		user, err = ur.UserDAO.GetByUsername(ctx, username)
	default:
		return entity.UserEntity{}, repoErr.ErrInvalidArgs
	}

	if err != nil {
		return entity.UserEntity{}, err
	}

	return user, nil
}

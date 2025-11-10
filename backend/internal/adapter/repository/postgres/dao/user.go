package dao

import (
	"context"

	"github.com/andreychh/coopera-backend/internal/adapter/repository/model/postgres"
	"github.com/andreychh/coopera-backend/pkg/logger"
)

type UserDAO struct {
	Logger *logger.Logger
}

func NewUserDAO(logger *logger.Logger) *UserDAO {
	return &UserDAO{
		Logger: logger,
	}
}

func (dao *UserDAO) SaveUser(ctx context.Context, user postgres.User) error {
	// sql запрос
	return nil
}

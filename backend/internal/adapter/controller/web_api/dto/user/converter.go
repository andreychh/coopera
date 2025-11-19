package user

import "github.com/andreychh/coopera-backend/internal/entity"

type HasTelegramID interface {
	GetTelegramID() int64
}

func ToEntity[T HasTelegramID](req T) *entity.UserEntity {
	return &entity.UserEntity{
		TelegramID: req.GetTelegramID(),
	}
}

package converter

import (
	"github.com/andreychh/coopera-backend/internal/adapter/repository/model/user_model"
	"github.com/andreychh/coopera-backend/internal/entity"
)

// FromEntityToModel преобразует entity в модель БД
func FromEntityToModel(euser entity.UserEntity) user_model.User {
	return user_model.User{
		TelegramID: euser.TelegramID,
	}
}

// FromModelToEntity преобразует модель БД в entity
func FromModelToEntity(muser user_model.User) entity.UserEntity {
	return entity.UserEntity{
		ID:         &muser.ID,
		TelegramID: muser.TelegramID,
		CreatedAt:  &muser.CreatedAt,
	}
}

func FromModelToEntityWithTeams(m user_model.UserWithTeams) entity.UserEntity {
	teams := make([]entity.TeamWithRole, len(m.Teams))
	for i, t := range m.Teams {
		teams[i] = entity.TeamWithRole{
			TeamID:   t.ID,
			TeamName: t.Name,
			Role:     entity.Role(t.Role),
		}
	}

	return entity.UserEntity{
		ID:         &m.ID,
		TelegramID: m.TelegramID,
		CreatedAt:  &m.CreatedAt,
		Teams:      teams,
	}
}

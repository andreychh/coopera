package converter

import (
	"github.com/andreychh/coopera-backend/internal/adapter/repository/model/user_model"
	"github.com/andreychh/coopera-backend/internal/entity"
)

func FromEntityToModel(euser entity.UserEntity) user_model.User {
	return user_model.User{
		TelegramID: *euser.TelegramID,
		Username:   *euser.Username,
	}
}

func FromModelToEntity(muser user_model.User) entity.UserEntity {
	return entity.UserEntity{
		ID:         &muser.ID,
		TelegramID: &muser.TelegramID,
		Username:   &muser.Username,
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
		TelegramID: &m.TelegramID,
		Username:   &m.Username,
		CreatedAt:  &m.CreatedAt,
		Teams:      teams,
	}
}

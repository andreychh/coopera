package converter

import (
	"github.com/andreychh/coopera-backend/internal/adapter/repository/model/team_model"
	"github.com/andreychh/coopera-backend/internal/entity"
)

func FromEntityToModelTeam(team entity.TeamEntity) team_model.Team {
	return team_model.Team{
		ID:        0,
		Name:      team.Name,
		CreatedBy: team.CreatedBy,
	}
}

func FromModelToEntityTeam(team team_model.Team) entity.TeamEntity {
	return entity.TeamEntity{
		ID:        &team.ID,
		Name:      team.Name,
		CreatedBy: team.CreatedBy,
		CreatedAt: &team.CreatedAt,
	}
}

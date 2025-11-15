package team_repo

import (
	"context"
	"github.com/andreychh/coopera/internal/adapter/repository/converter"
	"github.com/andreychh/coopera/internal/adapter/repository/postgres/dao"
	"github.com/andreychh/coopera/internal/entity"
)

type TeamRepository struct {
	dao dao.TeamDAO
}

func NewTeamRepository(dao dao.TeamDAO) *TeamRepository {
	return &TeamRepository{dao: dao}
}

func (r *TeamRepository) CreateRepo(ctx context.Context, e entity.TeamEntity) (entity.TeamEntity, error) {
	return r.dao.Create(ctx, converter.FromEntityToModelTeam(e))
}

func (r *TeamRepository) DeleteRepo(ctx context.Context, teamID int32) error {
	return r.dao.Delete(ctx, teamID)
}

func (r *TeamRepository) GetByIDRepo(ctx context.Context, teamID int32) (entity.TeamEntity, error) {
	return r.dao.GetByID(ctx, teamID)
}

func (r *TeamRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	return r.dao.ExistsByName(ctx, name)
}

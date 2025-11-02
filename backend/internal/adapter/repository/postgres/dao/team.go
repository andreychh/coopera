package dao

import (
	"context"
	"fmt"
	"github.com/andreychh/coopera/internal/adapter/repository/converter"
	team_model "github.com/andreychh/coopera/internal/adapter/repository/model/team_model"
	"github.com/andreychh/coopera/internal/adapter/repository/postgres"
	"github.com/andreychh/coopera/internal/entity"
)

type TeamDAO struct {
	db *postgres.DB
}

func NewTeamDAO(db *postgres.DB) *TeamDAO {
	return &TeamDAO{db: db}
}

func (r *TeamDAO) Create(ctx context.Context, t team_model.Team) (entity.TeamEntity, error) {
	const query = `
		INSERT INTO coopera.teams (name, created_by, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, name, created_by, created_at
	`
	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return entity.TeamEntity{}, fmt.Errorf("transaction not found")
	}

	var created team_model.Team
	if err := tx.QueryRow(ctx, query, t.Name, t.CreatedBy).Scan(
		&created.ID, &created.Name, &created.CreatedBy, &created.CreatedAt,
	); err != nil {
		return entity.TeamEntity{}, err
	}
	return converter.FromModelToEntityTeam(created), nil
}

func (r *TeamDAO) GetByID(ctx context.Context, teamID int32) (entity.TeamEntity, error) {
	const query = `
		SELECT id, name, created_by, created_at
		FROM coopera.teams
		WHERE id = $1
	`
	var t team_model.Team
	if err := r.db.Pool.QueryRow(ctx, query, teamID).Scan(&t.ID, &t.Name, &t.CreatedBy, &t.CreatedAt); err != nil {
		return entity.TeamEntity{}, err
	}
	return converter.FromModelToEntityTeam(t), nil
}

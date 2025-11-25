package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/andreychh/coopera-backend/internal/adapter/repository/converter"
	repoErr "github.com/andreychh/coopera-backend/internal/adapter/repository/errors"
	team_model "github.com/andreychh/coopera-backend/internal/adapter/repository/model/team_model"
	"github.com/andreychh/coopera-backend/internal/adapter/repository/postgres"
	"github.com/andreychh/coopera-backend/internal/entity"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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
		return entity.TeamEntity{}, repoErr.ErrTransactionNotFound
	}

	var created team_model.Team
	if err := tx.QueryRow(ctx, query, t.Name, t.CreatedBy).Scan(
		&created.ID, &created.Name, &created.CreatedBy, &created.CreatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return entity.TeamEntity{}, repoErr.ErrAlreadyExists
			case "23503": // foreign_key_violation
				return entity.TeamEntity{}, repoErr.ErrFailCreate
			}
		}
		return entity.TeamEntity{}, fmt.Errorf("%w: %v", repoErr.ErrFailCreate, err)
	}

	return converter.FromModelToEntityTeam(created), nil
}

func (r *TeamDAO) Delete(ctx context.Context, teamID int32) error {
	const query = `DELETE FROM coopera.teams WHERE id = $1`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return repoErr.ErrTransactionNotFound
	}

	if _, err := tx.Exec(ctx, query, teamID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // FK violation
			return repoErr.ErrFailDelete
		}
		return fmt.Errorf("%w: %v", repoErr.ErrFailDelete, err)
	}

	return nil
}

func (r *TeamDAO) GetByID(ctx context.Context, teamID int32) (entity.TeamEntity, error) {
	const query = `
		SELECT id, name, created_by, created_at
		FROM coopera.teams
		WHERE id = $1
	`

	var t team_model.Team
	err := r.db.Pool.QueryRow(ctx, query, teamID).Scan(&t.ID, &t.Name, &t.CreatedBy, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.TeamEntity{}, repoErr.ErrNotFound
		}
		return entity.TeamEntity{}, fmt.Errorf("%w: %v", repoErr.ErrDB, err)
	}

	return converter.FromModelToEntityTeam(t), nil
}

func (r *TeamDAO) ExistsByName(ctx context.Context, name string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM coopera.teams WHERE name = $1)`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return false, repoErr.ErrTransactionNotFound
	}

	var exists bool
	if err := tx.QueryRow(ctx, query, name).Scan(&exists); err != nil {
		return false, fmt.Errorf("%w: %v", repoErr.ErrFailCheckExists, err)
	}

	return exists, nil
}

func (r *TeamDAO) ExistsByID(ctx context.Context, teamID int32) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM coopera.teams WHERE id = $1)`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return false, repoErr.ErrTransactionNotFound
	}

	var exists bool
	if err := tx.QueryRow(ctx, query, teamID).Scan(&exists); err != nil {
		return false, fmt.Errorf("%w: %v", repoErr.ErrFailCheckExists, err)
	}

	return exists, nil
}

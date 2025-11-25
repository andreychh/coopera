package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/andreychh/coopera-backend/internal/adapter/repository/converter"
	repoErr "github.com/andreychh/coopera-backend/internal/adapter/repository/errors"
	"github.com/andreychh/coopera-backend/internal/adapter/repository/model/task_model"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/andreychh/coopera-backend/internal/adapter/repository/postgres"
	"github.com/andreychh/coopera-backend/internal/entity"
)

type TaskDAO struct {
	db *postgres.DB
}

func NewTaskDAO(db *postgres.DB) *TaskDAO {
	return &TaskDAO{db: db}
}

func (r *TaskDAO) Create(ctx context.Context, task task_model.Task) (entity.Task, error) {
	const query = `
		INSERT INTO coopera.tasks (team_id, title, description, points, assigned_to, created_by, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, team_id, title, description, points, status, assigned_to, created_by, created_at, updated_at
	`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return entity.Task{}, repoErr.ErrTransactionNotFound
	}

	var m task_model.Task
	err := tx.QueryRow(ctx, query, task.TeamID, task.Title,
		task.Description, task.Points, task.AssignedTo, task.CreatedBy, task.Status,
	).Scan(&m.ID, &m.TeamID, &m.Title, &m.Description, &m.Points,
		&m.Status, &m.AssignedTo, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return entity.Task{}, repoErr.ErrAlreadyExists
			}
		}
		return entity.Task{}, fmt.Errorf("%w: %v", repoErr.ErrFailCreate, err)
	}

	return converter.FromModelToEntityTask(m), nil
}

func (r *TaskDAO) GetByAssignedID(ctx context.Context, userID int32) ([]entity.Task, error) {
	const query = `
		SELECT id, team_id, title, description, points, status, assigned_to, 
		       created_by, created_at, updated_at
		FROM coopera.tasks
		WHERE assigned_to = $1
		ORDER BY created_at DESC
	`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return nil, repoErr.ErrTransactionNotFound
	}

	rows, err := tx.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repoErr.ErrFailGet, err)
	}
	defer rows.Close()

	var tasks []entity.Task

	for rows.Next() {
		var m task_model.Task

		if err := rows.Scan(
			&m.ID, &m.TeamID, &m.Title, &m.Description, &m.Points,
			&m.Status, &m.AssignedTo, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", repoErr.ErrFailGet, err)
		}

		tasks = append(tasks, converter.FromModelToEntityTask(m))
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%w: %v", repoErr.ErrFailGet, rows.Err())
	}

	return tasks, nil
}

func (r *TaskDAO) GetByTaskID(ctx context.Context, id int32) (entity.Task, error) {
	const query = `
		SELECT id, team_id, title, description, points, status, assigned_to,
		       created_by, created_at, updated_at
		FROM coopera.tasks
		WHERE id = $1
	`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return entity.Task{}, repoErr.ErrTransactionNotFound
	}

	var m task_model.Task
	err := tx.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.TeamID, &m.Title, &m.Description, &m.Points,
		&m.Status, &m.AssignedTo, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Task{}, repoErr.ErrNotFound
		}

		return entity.Task{}, fmt.Errorf("%w: %v", repoErr.ErrFailGet, err)
	}

	return converter.FromModelToEntityTask(m), nil
}

func (r *TaskDAO) GetByTeamID(ctx context.Context, teamID int32) ([]entity.Task, error) {
	const query = `
		SELECT id, team_id, title, description, points, status, assigned_to,
		       created_by, created_at, updated_at
		FROM coopera.tasks
		WHERE team_id = $1
		ORDER BY created_at DESC
	`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return nil, repoErr.ErrTransactionNotFound
	}

	rows, err := tx.Query(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repoErr.ErrFailGet, err)
	}
	defer rows.Close()

	var tasks []entity.Task

	for rows.Next() {
		var m task_model.Task

		if err := rows.Scan(
			&m.ID, &m.TeamID, &m.Title, &m.Description, &m.Points,
			&m.Status, &m.AssignedTo, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", repoErr.ErrFailGet, err)
		}

		tasks = append(tasks, converter.FromModelToEntityTask(m))
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%w: %v", repoErr.ErrFailGet, rows.Err())
	}

	return tasks, nil
}

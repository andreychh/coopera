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
	"strings"

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
		INSERT INTO coopera.tasks (team_id, title, description, points, assigned_to, created_by, status, created_by_member_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, team_id, title, description, points, status, assigned_to, created_by, created_at, updated_at, created_by_member_id
	`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return entity.Task{}, repoErr.ErrTransactionNotFound
	}

	var m task_model.Task
	err := tx.QueryRow(ctx, query, task.TeamID, task.Title,
		task.Description, task.Points, task.AssignedTo, task.CreatedByUser, task.Status, task.CreatedByMember,
	).Scan(&m.ID, &m.TeamID, &m.Title, &m.Description, &m.Points,
		&m.Status, &m.AssignedTo, &m.CreatedByUser, &m.CreatedAt, &m.UpdatedAt, &m.CreatedByMember,
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

func (r *TaskDAO) Update(ctx context.Context, task task_model.UpdateTask) error {
	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return repoErr.ErrTransactionNotFound
	}

	var setParts []string
	var args []interface{}
	argIdx := 1

	if task.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIdx))
		args = append(args, *task.Title)
		argIdx++
	}

	if task.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, *task.Description)
		argIdx++
	}

	if task.Points != nil {
		setParts = append(setParts, fmt.Sprintf("points = $%d", argIdx))
		args = append(args, *task.Points)
		argIdx++
	}

	if task.AssignedTo != nil {
		setParts = append(setParts, fmt.Sprintf("assigned_to = $%d", argIdx))
		args = append(args, *task.AssignedTo)
		argIdx++
	}

	if len(setParts) == 0 {
		return repoErr.ErrNothingToUpdate
	}

	setParts = append(setParts, "updated_at = NOW()")

	query := fmt.Sprintf(`
		UPDATE coopera.tasks
		SET %s
		WHERE id = $%d
	`, strings.Join(setParts, ", "), argIdx)

	args = append(args, task.ID)

	// Используем Exec, но обрабатываем ошибку
	_, err := tx.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation (team_id, title)
				return repoErr.ErrAlreadyExists
			case "23503": // foreign_key_violation
				if strings.Contains(pgErr.ConstraintName, "fk_assigned_to_membership") || strings.Contains(pgErr.ConstraintName, "fk_team") {
					return repoErr.ErrInvalidArgs
				}
			case "23514": // check_violation
				return fmt.Errorf("%w: %s", repoErr.ErrInvalidArgs, pgErr.Message)
			}
		}
		return fmt.Errorf("%w: %v", repoErr.ErrFailUpdate, err)
	}

	return nil
}

func (r *TaskDAO) GetByAssignedTo(ctx context.Context, memberID int32) ([]entity.Task, error) {
	const query = `
		SELECT id, team_id, title, description, points, status, assigned_to, 
		       created_by, created_at, updated_at, created_by_member_id
		FROM coopera.tasks
		WHERE assigned_to = $1
		ORDER BY created_at DESC
	`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return nil, repoErr.ErrTransactionNotFound
	}

	rows, err := tx.Query(ctx, query, memberID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repoErr.ErrFailGet, err)
	}
	defer rows.Close()

	var tasks []entity.Task

	for rows.Next() {
		var m task_model.Task

		if err := rows.Scan(
			&m.ID, &m.TeamID, &m.Title, &m.Description, &m.Points,
			&m.Status, &m.AssignedTo, &m.CreatedByUser, &m.CreatedAt, &m.UpdatedAt, &m.CreatedByMember,
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
		       created_by, created_at, updated_at, created_by_member_id
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
		&m.Status, &m.AssignedTo, &m.CreatedByUser, &m.CreatedAt, &m.UpdatedAt, &m.CreatedByMember,
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
		       created_by, created_at, updated_at, created_by_member_id
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
			&m.Status, &m.AssignedTo, &m.CreatedByUser, &m.CreatedAt, &m.UpdatedAt, &m.CreatedByMember,
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

func (r *TaskDAO) UpdateStatus(ctx context.Context, status task_model.TaskStatus) error {
	const query = `
		UPDATE coopera.tasks
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return repoErr.ErrTransactionNotFound
	}

	cmdTag, err := tx.Exec(ctx, query, status.Status, status.TaskID)
	if err != nil {
		return fmt.Errorf("%w: %v", repoErr.ErrFailUpdate, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return repoErr.ErrNotFound
	}

	return nil
}

func (r *TaskDAO) Delete(ctx context.Context, taskID int32) error {
	const query = `
		DELETE FROM coopera.tasks
		WHERE id = $1
	`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return repoErr.ErrTransactionNotFound
	}

	_, err := tx.Exec(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("%w: %v", repoErr.ErrFailDelete, err)
	}

	return nil
}

func (r *TaskDAO) GetAllTasks(ctx context.Context) ([]entity.Task, error) {
	const query = `
		SELECT id, team_id, title, description, points, status, assigned_to,
		       created_by, created_at, updated_at, created_by_member_id
		FROM coopera.tasks
		ORDER BY created_at ASC
	`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return nil, repoErr.ErrTransactionNotFound
	}

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repoErr.ErrFailGet, err)
	}
	defer rows.Close()

	var tasks []entity.Task

	for rows.Next() {
		var m task_model.Task

		if err := rows.Scan(
			&m.ID, &m.TeamID, &m.Title, &m.Description, &m.Points,
			&m.Status, &m.AssignedTo, &m.CreatedByUser, &m.CreatedAt, &m.UpdatedAt, &m.CreatedByMember,
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

package dao

import (
	"context"
	"fmt"
	"github.com/andreychh/coopera/internal/adapter/repository/converter"
	"github.com/andreychh/coopera/internal/adapter/repository/model/user_model"
	"github.com/andreychh/coopera/internal/adapter/repository/postgres"
	"github.com/andreychh/coopera/internal/entity"
)

type UserRepository struct {
	db *postgres.DB
}

func NewUserDAO(db *postgres.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Create(ctx context.Context, muser user_model.User) (entity.UserEntity, error) {
	const query = `
		INSERT INTO coopera.users (telegram_id, created_at) 
		VALUES ($1, NOW()) 
		RETURNING id, telegram_id, created_at
	`

	// Всегда работаем через транзакцию из контекста
	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return entity.UserEntity{}, fmt.Errorf("transaction not found in context")
	}

	var userModel user_model.User
	err := tx.QueryRow(ctx, query, muser.TelegramID).Scan(
		&userModel.ID,
		&userModel.TelegramID,
		&userModel.CreatedAt,
	)

	if err != nil {
		return entity.UserEntity{}, fmt.Errorf("failed to create user_model: %w", err)
	}

	return converter.FromModelToEntity(userModel), nil
}

func (ur *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (entity.UserEntity, error) {
	const query = `
		SELECT 
			u.id, u.telegram_id, u.created_at,
			t.id AS team_id, t.name AS team_name, m.role
		FROM coopera.users u
		LEFT JOIN coopera.memberships m ON m.member_id = u.id
		LEFT JOIN coopera.teams t ON t.id = m.team_id
		WHERE u.telegram_id = $1
	`

	// тоже извлекаем транзакцию из контекста
	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return entity.UserEntity{}, fmt.Errorf("transaction not found in context")
	}

	rows, err := tx.Query(ctx, query, telegramID)
	if err != nil {
		return entity.UserEntity{}, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var user user_model.UserWithTeams
	for rows.Next() {
		var (
			teamID   *int32
			teamName *string
			role     *string
		)

		if err := rows.Scan(
			&user.ID,
			&user.TelegramID,
			&user.CreatedAt,
			&teamID,
			&teamName,
			&role,
		); err != nil {
			return entity.UserEntity{}, err
		}

		if teamID != nil {
			user.Teams = append(user.Teams, user_model.TeamWithRole{
				ID:   *teamID,
				Name: *teamName,
				Role: *role,
			})
		}
	}

	if user.ID == 0 {
		return entity.UserEntity{}, fmt.Errorf("user not found")
	}

	return converter.FromModelToEntityWithTeams(user), nil
}

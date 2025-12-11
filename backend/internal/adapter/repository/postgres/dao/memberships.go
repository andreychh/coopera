package dao

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"

	"github.com/andreychh/coopera-backend/internal/adapter/repository/converter"
	repoErr "github.com/andreychh/coopera-backend/internal/adapter/repository/errors"
	"github.com/andreychh/coopera-backend/internal/adapter/repository/model/membership_model"
	"github.com/andreychh/coopera-backend/internal/adapter/repository/postgres"
	"github.com/andreychh/coopera-backend/internal/entity"
)

type MembershipDAO struct {
	db *postgres.DB
}

func NewMembershipDAO(db *postgres.DB) *MembershipDAO {
	return &MembershipDAO{db: db}
}

func (r *MembershipDAO) AddMember(ctx context.Context, m membership_model.Membership) (int32, error) {
	const query = `
        INSERT INTO coopera.memberships (team_id, user_id, role, created_at) 
        VALUES ($1, $2, $3, NOW()) 
        ON CONFLICT (team_id, user_id) DO NOTHING
        RETURNING id
    `

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return 0, repoErr.ErrTransactionNotFound
	}

	var id int32
	if err := tx.QueryRow(ctx, query, m.TeamID, m.UserID, m.Role).Scan(&id); err != nil {
		if err == pgx.ErrNoRows {
			return 0, repoErr.ErrMemberAlreadyExists
		}
		return 0, fmt.Errorf("%w: %v", repoErr.ErrFailToAdd, err)
	}

	return id, nil
}

func (r *MembershipDAO) DeleteMember(ctx context.Context, m membership_model.Membership) error {
	const query = `
		DELETE FROM coopera.memberships
		WHERE team_id = $1 AND user_id = $2
	`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return repoErr.ErrTransactionNotFound
	}

	result, err := tx.Exec(ctx, query, m.TeamID, m.UserID)
	if err != nil {
		return fmt.Errorf("%w: %v", repoErr.ErrFailDelete, err)
	}

	if result.RowsAffected() == 0 {
		return repoErr.ErrNotFound
	}

	return nil
}

func (r *MembershipDAO) GetMembers(ctx context.Context, teamID int32) ([]entity.MembershipEntity, error) {
	const query = `
		SELECT id, team_id, user_id, role, created_at
		FROM coopera.memberships
		WHERE team_id = $1
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

	var members []entity.MembershipEntity
	for rows.Next() {
		var m membership_model.Membership
		if err := rows.Scan(&m.ID, &m.TeamID, &m.UserID, &m.Role, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", repoErr.ErrFailToCastScan, err)
		}
		members = append(members, converter.FromModelToEntityMembership(m))
	}

	if len(members) == 0 {
		return nil, repoErr.ErrNotFound
	}

	return members, nil
}

func (r *MembershipDAO) ExistsMember(ctx context.Context, memberID int32) (bool, error) {
	const query = `
		SELECT 1
		FROM coopera.memberships
		WHERE id = $1
		LIMIT 1
	`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return false, repoErr.ErrTransactionNotFound
	}

	row := tx.QueryRow(ctx, query, memberID)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if err.Error() == "no rows in result set" {
			return false, nil
		}
		return false, fmt.Errorf("%w: %v", repoErr.ErrFailGet, err)
	}

	return true, nil
}

func (r *MembershipDAO) GetMember(ctx context.Context, teamID, memberID int32) (entity.MembershipEntity, error) {
	const query = `
		SELECT id, team_id, user_id, role, created_at
		FROM coopera.memberships
		WHERE team_id = $1 AND user_id = $2
	`

	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return entity.MembershipEntity{}, repoErr.ErrTransactionNotFound
	}

	var m membership_model.Membership
	err := tx.QueryRow(ctx, query, teamID, memberID).Scan(&m.ID, &m.TeamID, &m.UserID, &m.Role, &m.CreatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return entity.MembershipEntity{}, repoErr.ErrNotFound
		}
		return entity.MembershipEntity{}, fmt.Errorf("%w: %v", repoErr.ErrFailGet, err)
	}

	return converter.FromModelToEntityMembership(m), nil
}

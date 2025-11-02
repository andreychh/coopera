package dao

import (
	"context"
	"fmt"
	"github.com/andreychh/coopera/internal/adapter/repository/converter"
	"github.com/andreychh/coopera/internal/adapter/repository/model/membership_model"
	"github.com/andreychh/coopera/internal/adapter/repository/postgres"
	"github.com/andreychh/coopera/internal/entity"
)

type MembershipDAO struct {
	db *postgres.DB
}

func NewMembershipDAO(db *postgres.DB) *MembershipDAO {
	return &MembershipDAO{db: db}
}

func (r *MembershipDAO) AddMember(ctx context.Context, m membership_model.Membership) error {
	const query = `
		INSERT INTO coopera.memberships (team_id, member_id, role, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (team_id, member_id) DO NOTHING
	`
	tx, ok := ctx.Value(postgres.TransactionKey{}).(postgres.Transaction)
	if !ok {
		return fmt.Errorf("transaction not found")
	}
	_, err := tx.Exec(ctx, query, m.TeamID, m.MemberID, m.Role)
	return err
}

func (r *MembershipDAO) GetMembers(ctx context.Context, teamID int32) ([]entity.MembershipEntity, error) {
	const query = `
		SELECT id, team_id, member_id, role, created_at
		FROM coopera.memberships
		WHERE team_id = $1
	`
	rows, err := r.db.Pool.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []entity.MembershipEntity
	for rows.Next() {
		var m membership_model.Membership
		if err := rows.Scan(&m.ID, &m.TeamID, &m.MemberID, &m.Role, &m.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, converter.FromModelToEntityMembership(m))
	}
	return members, nil
}

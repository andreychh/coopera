package converter

import (
	"github.com/andreychh/coopera-backend/internal/adapter/repository/model/membership_model"
	"github.com/andreychh/coopera-backend/internal/entity"
)

func FromEntityToModelMembership(membership entity.MembershipEntity) membership_model.Membership {
	return membership_model.Membership{
		ID:     0,
		TeamID: membership.TeamID,
		UserID: membership.UserID,
		Role:   string(membership.Role),
	}
}

func FromModelToEntityMembership(m membership_model.Membership) entity.MembershipEntity {
	return entity.MembershipEntity{
		ID:        m.ID,
		TeamID:    m.TeamID,
		UserID:    m.UserID,
		Role:      entity.Role(m.Role),
		CreatedAt: &m.CreatedAt,
	}
}

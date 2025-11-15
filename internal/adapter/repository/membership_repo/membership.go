package membership_repo

import (
	"context"
	"github.com/andreychh/coopera/internal/adapter/repository/converter"
	"github.com/andreychh/coopera/internal/adapter/repository/postgres/dao"
	"github.com/andreychh/coopera/internal/entity"
)

type MembershipRepository struct {
	dao dao.MembershipDAO
}

func NewMembershipRepository(dao dao.MembershipDAO) *MembershipRepository {
	return &MembershipRepository{dao: dao}
}

func (r *MembershipRepository) AddMemberRepo(ctx context.Context, m entity.MembershipEntity) error {
	return r.dao.AddMember(ctx, converter.FromEntityToModelMembership(m))
}

func (r *MembershipRepository) GetMembersRepo(ctx context.Context, teamID int32) ([]entity.MembershipEntity, error) {
	return r.dao.GetMembers(ctx, teamID)
}

func (r *MembershipRepository) DeleteMemberRepo(ctx context.Context, membership entity.MembershipEntity) error {
	return r.dao.DeleteMember(ctx, converter.FromEntityToModelMembership(membership))
}

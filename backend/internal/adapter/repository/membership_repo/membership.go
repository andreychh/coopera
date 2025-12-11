package membership_repo

import (
	"context"
	"github.com/andreychh/coopera-backend/internal/adapter/repository/converter"
	"github.com/andreychh/coopera-backend/internal/adapter/repository/postgres/dao"
	"github.com/andreychh/coopera-backend/internal/entity"
)

type MembershipRepository struct {
	dao dao.MembershipDAO
}

func NewMembershipRepository(dao dao.MembershipDAO) *MembershipRepository {
	return &MembershipRepository{dao: dao}
}

func (r *MembershipRepository) AddMemberRepo(ctx context.Context, m entity.MembershipEntity) (int32, error) {
	return r.dao.AddMember(ctx, converter.FromEntityToModelMembership(m))
}

func (r *MembershipRepository) GetMembersRepo(ctx context.Context, teamID int32) ([]entity.MembershipEntity, error) {
	return r.dao.GetMembers(ctx, teamID)
}

func (r *MembershipRepository) DeleteMemberRepo(ctx context.Context, membership entity.MembershipEntity) error {
	return r.dao.DeleteMember(ctx, converter.FromEntityToModelMembership(membership))
}

func (r *MembershipRepository) MemberExistsRepo(ctx context.Context, memberID int32) (bool, error) {
	return r.dao.ExistsMember(ctx, memberID)
}

func (r *MembershipRepository) GetMemberRepo(ctx context.Context, teamID, memberID int32) (entity.MembershipEntity, error) {
	return r.dao.GetMember(ctx, teamID, memberID)
}

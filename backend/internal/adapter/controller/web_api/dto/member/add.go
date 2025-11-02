package member

import (
	"github.com/andreychh/coopera/internal/entity"
)

type AddMemberRequest struct {
	TeamID   int32 `json:"team_id" validate:"required"`
	MemberID int32 `json:"member_id" validate:"required"`
}

type AddMemberResponse struct {
	ID        int32       `json:"id"`
	TeamID    int32       `json:"team_id"`
	MemberID  int32       `json:"member_id"`
	Role      entity.Role `json:"role"`
	CreatedAt string      `json:"created_at"`
}

func ToEntityAddMembersRequest(req *AddMemberRequest) *entity.MembershipEntity {
	return &entity.MembershipEntity{
		TeamID:   req.TeamID,
		MemberID: req.MemberID,
	}
}

func ToAddMembersResponse(member *entity.MembershipEntity) *AddMemberResponse {
	return &AddMemberResponse{
		ID:        member.ID,
		TeamID:    member.TeamID,
		MemberID:  member.MemberID,
		Role:      member.Role,
		CreatedAt: member.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

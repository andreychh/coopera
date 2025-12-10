package member

import (
	"github.com/andreychh/coopera-backend/internal/entity"
)

type AddMemberRequest struct {
	TeamID   int32 `json:"team_id" validate:"required"`
	MemberID int32 `json:"user_id" validate:"required"`
}

type AddMemberResponse struct {
	ID        int32       `json:"id"`
	TeamID    int32       `json:"team_id"`
	UserID    int32       `json:"user_id"`
	Role      entity.Role `json:"role"`
	CreatedAt string      `json:"created_at"`
}

func ToEntityAddMembersRequest(req *AddMemberRequest) *entity.MembershipEntity {
	return &entity.MembershipEntity{
		TeamID: req.TeamID,
		UserID: req.MemberID,
	}
}

func ToAddMembersResponse(member *entity.MembershipEntity) *AddMemberResponse {
	return &AddMemberResponse{
		ID:        member.ID,
		TeamID:    member.TeamID,
		UserID:    member.UserID,
		Role:      member.Role,
		CreatedAt: member.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

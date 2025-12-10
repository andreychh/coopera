package member

import "github.com/andreychh/coopera-backend/internal/entity"

type DeleteMemberRequest struct {
	TeamID        int32 `form:"team_id" validate:"required"`
	MemberID      int32 `form:"user_id" validate:"required"`
	CurrentUserID int32 `form:"current_user_id" validate:"required"`
}

func ToEntityDeleteMemberRequest(req *DeleteMemberRequest) *entity.MembershipEntity {
	return &entity.MembershipEntity{
		TeamID: req.TeamID,
		UserID: req.MemberID,
	}
}

package team

import "github.com/andreychh/coopera-backend/internal/entity"

type GetTeamRequest struct {
	TeamID int32 `form:"team_id" validate:"required"`
}

type GetTeamResponse struct {
	ID        int32            `json:"id"`
	Name      string           `json:"name"`
	CreatedAt string           `json:"created_at"`
	CreatedBy int32            `json:"created_by"`
	Members   []TeamMemberInfo `json:"members"`
}

type TeamMemberInfo struct {
	MemberID int32       `json:"member_id"`
	Role     entity.Role `json:"role"`
}

func ToGetTeamResponse(team entity.TeamEntity, members []entity.MembershipEntity) GetTeamResponse {
	membersInfo := make([]TeamMemberInfo, len(members))
	for i, m := range members {
		membersInfo[i] = TeamMemberInfo{
			MemberID: m.MemberID,
			Role:     m.Role,
		}
	}

	return GetTeamResponse{
		ID:        *team.ID,
		Name:      team.Name,
		CreatedAt: team.CreatedAt.Format("2006-01-02T15:04:05Z"),
		CreatedBy: team.CreatedBy,
		Members:   membersInfo,
	}
}

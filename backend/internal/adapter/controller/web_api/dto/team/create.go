package team

import "github.com/andreychh/coopera/internal/entity"

type CreateTeamRequest struct {
	UserID int32  `json:"user_id" validate:"required"`
	Name   string `json:"name" validate:"required"`
}

type CreateTeamResponse struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	CreatedBy int32  `json:"created_by"`
}

func ToEntityCreateTeamRequest(req *CreateTeamRequest) *entity.TeamEntity {
	return &entity.TeamEntity{
		Name:      req.Name,
		CreatedBy: req.UserID,
	}
}

func ToCreateTeamResponse(team *entity.TeamEntity) *CreateTeamResponse {
	return &CreateTeamResponse{
		ID:        *team.ID,
		Name:      team.Name,
		CreatedAt: team.CreatedAt.Format("2006-01-02T15:04:05Z"),
		CreatedBy: team.CreatedBy,
	}
}

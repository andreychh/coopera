package team

type DeleteTeamRequest struct {
	TeamID        int32 `form:"team_id" validate:"required"`
	CurrentUserID int32 `form:"current_user_id" validate:"required"`
}

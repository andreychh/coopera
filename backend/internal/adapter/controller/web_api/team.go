package web_api

import (
	"net/http"

	teamdto "github.com/andreychh/coopera/internal/adapter/controller/web_api/dto/team"
	"github.com/andreychh/coopera/internal/usecase"
)

type TeamController struct {
	teamUseCase usecase.TeamUseCase
}

func NewTeamController(teamUseCase usecase.TeamUseCase) *TeamController {
	return &TeamController{
		teamUseCase: teamUseCase,
	}
}

func (tc *TeamController) Create(w http.ResponseWriter, r *http.Request) {
	var req teamdto.CreateTeamRequest
	if err := BindRequest(r, &req); err != nil {
		writeValidationError(w, err)
		return
	}

	team, err := tc.teamUseCase.CreateUsecase(r.Context(), *teamdto.ToEntityCreateTeamRequest(&req))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, teamdto.ToCreateTeamResponse(&team))
}

func (tc *TeamController) Get(w http.ResponseWriter, r *http.Request) {
	var req teamdto.GetTeamRequest

	if err := BindRequest(r, &req); err != nil {
		writeValidationError(w, err)
		return
	}

	team, membership, err := tc.teamUseCase.GetByIDUsecase(r.Context(), req.TeamID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, teamdto.ToGetTeamResponse(team, membership))
}

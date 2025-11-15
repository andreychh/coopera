package web_api

import (
	errors "github.com/andreychh/coopera/pkg/errors"
	"github.com/go-playground/validator/v10"
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

func (tc *TeamController) Create(w http.ResponseWriter, r *http.Request) error {
	var req teamdto.CreateTeamRequest
	if err := BindRequest(r, &req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errors.WrapValidationError(ve)
		}
		return errors.ErrInvalidInput
	}

	team, err := tc.teamUseCase.CreateUsecase(r.Context(), *teamdto.ToEntityCreateTeamRequest(&req))
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusCreated, teamdto.ToCreateTeamResponse(&team))
	return nil
}

func (tc *TeamController) Get(w http.ResponseWriter, r *http.Request) error {
	var req teamdto.GetTeamRequest
	if err := BindRequest(r, &req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errors.WrapValidationError(ve)
		}
		return errors.ErrInvalidInput
	}

	team, membership, err := tc.teamUseCase.GetByIDUsecase(r.Context(), req.TeamID)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, teamdto.ToGetTeamResponse(team, membership))
	return nil
}

func (tc *TeamController) Delete(w http.ResponseWriter, r *http.Request) error {
	var req teamdto.DeleteTeamRequest
	if err := BindRequest(r, &req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errors.WrapValidationError(ve)
		}
		return errors.ErrInvalidInput
	}

	if err := tc.teamUseCase.DeleteUsecase(r.Context(), req.TeamID, req.CurrentUserID); err != nil {
		return err
	}

	writeJSON(w, http.StatusNoContent, nil)
	return nil
}

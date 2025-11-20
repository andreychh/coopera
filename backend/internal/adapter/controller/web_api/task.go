package web_api

import (
	errors "github.com/andreychh/coopera-backend/pkg/errors"
	"github.com/go-playground/validator/v10"
	"net/http"

	taskdto "github.com/andreychh/coopera-backend/internal/adapter/controller/web_api/dto/task"
	"github.com/andreychh/coopera-backend/internal/usecase"
)

type TaskController struct {
	taskUseCase usecase.TaskUseCase
}

func NewTaskController(taskUseCase usecase.TaskUseCase) *TaskController {
	return &TaskController{
		taskUseCase: taskUseCase,
	}
}

func (tc *TaskController) Create(w http.ResponseWriter, r *http.Request) error {
	var req taskdto.CreateTaskRequest
	if err := BindRequest(r, &req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errors.WrapValidationError(ve)
		}
		return errors.ErrInvalidInput
	}

	team, err := tc.taskUseCase.CreateUsecase(r.Context(), *taskdto.ToEntityCreateTaskRequest(&req))
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusCreated, taskdto.ToCreateTaskResponse(&team))
	return nil
}

//func (tc *TaskController) Get(w http.ResponseWriter, r *http.Request) error {
//	var req teamdto.GetTeamRequest
//	if err := BindRequest(r, &req); err != nil {
//		if ve, ok := err.(validator.ValidationErrors); ok {
//			return errors.WrapValidationError(ve)
//		}
//		return errors.ErrInvalidInput
//	}
//
//	team, membership, err := tc.teamUseCase.GetByIDUsecase(r.Context(), req.TeamID)
//	if err != nil {
//		return err
//	}
//
//	writeJSON(w, http.StatusOK, teamdto.ToGetTeamResponse(team, membership))
//	return nil
//}
//
//func (tc *TaskController) Delete(w http.ResponseWriter, r *http.Request) error {
//	var req teamdto.DeleteTeamRequest
//	if err := BindRequest(r, &req); err != nil {
//		if ve, ok := err.(validator.ValidationErrors); ok {
//			return errors.WrapValidationError(ve)
//		}
//		return errors.ErrInvalidInput
//	}
//
//	if err := tc.teamUseCase.DeleteUsecase(r.Context(), req.TeamID, req.CurrentUserID); err != nil {
//		return err
//	}
//
//	writeJSON(w, http.StatusNoContent, nil)
//	return nil
//}

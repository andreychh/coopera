package web_api

import (
	"fmt"
	"github.com/andreychh/coopera-backend/pkg/errors"
	"github.com/go-playground/validator/v10"
	"net/http"

	userdto "github.com/andreychh/coopera-backend/internal/adapter/controller/web_api/dto/user"
	"github.com/andreychh/coopera-backend/internal/usecase"
)

type UserController struct {
	userUseCase usecase.UserUseCase
}

func NewUserController(userUseCase usecase.UserUseCase) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) error {
	var req userdto.CreateUserRequest
	if err := BindRequest(r, &req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errors.WrapValidationError(ve)
		}
		return errors.ErrInvalidInput
	}

	user, err := uc.userUseCase.CreateUsecase(r.Context(), *userdto.FromCreateUserRequest(&req))
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusCreated, userdto.ToCreateUserResponse(&user))
	return nil
}

func (uc *UserController) Get(w http.ResponseWriter, r *http.Request) error {
	var req userdto.GetUserRequest
	if err := BindRequest(r, &req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errors.WrapValidationError(ve)
		}
		return errors.ErrInvalidInput
	}

	if err := req.Validate(); err != nil {
		return fmt.Errorf("%w: %v", errors.ErrInvalidInput, err)
	}

	user, err := uc.userUseCase.GetUsecase(r.Context(), req.TelegramID, req.UserName, req.ID)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, userdto.ToGetUserResponse(&user))
	return nil
}

func (uc *UserController) Delete(w http.ResponseWriter, r *http.Request) error {
	var req userdto.DeleteUserRequest
	if err := BindRequest(r, &req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errors.WrapValidationError(ve)
		}
		return errors.ErrInvalidInput
	}

	if err := uc.userUseCase.DeleteUsecase(r.Context(), req.ID); err != nil {
		return err
	}

	writeJSON(w, http.StatusNoContent, nil)
	return nil
}

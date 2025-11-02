package web_api

import (
	"net/http"

	userdto "github.com/andreychh/coopera/internal/adapter/controller/web_api/dto/user"
	"github.com/andreychh/coopera/internal/usecase"
)

type UserController struct {
	userUseCase usecase.UserUseCase
}

func NewUserController(userUseCase usecase.UserUseCase) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

// Create — создаёт нового пользователя
func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var req userdto.CreateUserRequest

	if err := BindRequest(r, &req); err != nil {
		writeValidationError(w, err)
		return
	}

	user, err := uc.userUseCase.CreateUsecase(r.Context(), *userdto.ToEntity(&req))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, userdto.ToCreateUserResponse(&user))
}

func (uc *UserController) Get(w http.ResponseWriter, r *http.Request) {
	var req userdto.GetUserRequest
	if err := BindRequest(r, &req); err != nil {
		writeValidationError(w, err)
		return
	}

	user, err := uc.userUseCase.GetUsecase(r.Context(), *userdto.ToEntity(&req))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, userdto.ToGetUserResponse(&user))
}

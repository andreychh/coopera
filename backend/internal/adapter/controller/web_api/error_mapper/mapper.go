package error_mapper

import (
	"errors"
	"net/http"

	dbErr "github.com/andreychh/coopera-backend/internal/adapter/repository/errors"
	appErr "github.com/andreychh/coopera-backend/pkg/errors"
)

func MapErrorToHTTP(err error) (int, string) {
	switch {
	// Ошибки юзкейсов
	case errors.Is(err, appErr.ErrNotFound) ||
		errors.Is(err, appErr.ErrMemberNotFound) ||
		errors.Is(err, appErr.ErrTeamNotFound):
		return http.StatusUnprocessableEntity, err.Error()
	case errors.Is(err, appErr.ErrAlreadyExists):
		return http.StatusConflict, err.Error()
	case errors.Is(err, appErr.ErrInvalidInput) ||
		errors.Is(err, appErr.ErrTaskFilter) ||
		errors.Is(err, appErr.ErrAssignedMemberNotExists):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, appErr.ErrUnauthorized):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, appErr.ErrForbidden) ||
		errors.Is(err, appErr.ErrUserOwner) ||
		errors.Is(err, appErr.ErrNoPermissionToDelete) ||
		errors.Is(err, appErr.ErrNoPermissionToUpdate) ||
		errors.Is(err, appErr.ErrOnlyManagerCanUpdatePoints) ||
		errors.Is(err, appErr.ErrOnlyManagerCanAssign) ||
		errors.Is(err, appErr.ErrOnlyManagerCanSetPoints) ||
		errors.Is(err, appErr.ErrCantAssignWithoutPoints) ||
		errors.Is(err, appErr.ErrOnlyManagerOrSelfCanAssign):
		return http.StatusForbidden, err.Error()
	case errors.Is(err, appErr.ErrConflict):
		return http.StatusConflict, err.Error()

	// Ошибки репозитория / БД
	case errors.Is(err, dbErr.ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, dbErr.ErrInvalidArgs):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, dbErr.ErrAlreadyExists) ||
		errors.Is(err, dbErr.ErrMemberAlreadyExists):
		return http.StatusConflict, err.Error()
	case errors.Is(err, dbErr.ErrFailCreate) || errors.Is(err, dbErr.ErrFailDelete) ||
		errors.Is(err, dbErr.ErrFailToAdd) || errors.Is(err, dbErr.ErrFailGet) ||
		errors.Is(err, dbErr.ErrFailCheckExists) || errors.Is(err, dbErr.ErrFailToCastScan):
		return http.StatusInternalServerError, err.Error()
	case errors.Is(err, dbErr.ErrTransactionNotFound):
		return http.StatusInternalServerError, err.Error()
	case errors.Is(err, dbErr.ErrDB):
		return http.StatusInternalServerError, "database error"

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

package errors

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

// Сентинелы
var (
	ErrNotFound             = errors.New("not found")
	ErrAlreadyExists        = errors.New("already exists")
	ErrInvalidInput         = errors.New("invalid input")
	ErrForbidden            = errors.New("forbidden")
	ErrConflict             = errors.New("conflict")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrUserOwner            = errors.New("cannot delete user is the owner of the team")
	ErrNoPermissionToDelete = errors.New("no permission to delete")
	ErrNoPermissionToUpdate = errors.New("no permission to update")
	ErrTaskFilter           = errors.New("no filter provided")
	ErrMemberNotFound       = errors.New("member not found")
	ErrTeamNotFound         = errors.New("team not found")
)

type ValidationError struct {
	Base    error
	Details []string
}

func (e *ValidationError) Error() string {
	if len(e.Details) == 0 {
		return e.Base.Error()
	}
	return fmt.Sprintf("%s: %v", e.Base.Error(), e.Details)
}

func (e *ValidationError) Unwrap() error {
	return e.Base
}

// конвертация validator.ValidationErrors в ValidationError
func WrapValidationError(ve validator.ValidationErrors) error {
	details := make([]string, 0, len(ve))
	for _, f := range ve {
		details = append(details, f.Error())
	}
	return &ValidationError{
		Base:    ErrInvalidInput,
		Details: details,
	}
}

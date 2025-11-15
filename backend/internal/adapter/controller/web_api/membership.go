package web_api

import (
	"github.com/andreychh/coopera/pkg/errors"
	"github.com/go-playground/validator/v10"
	"net/http"

	memberdto "github.com/andreychh/coopera/internal/adapter/controller/web_api/dto/member"
	"github.com/andreychh/coopera/internal/usecase"
)

type MembershipController struct {
	membershipUseCase usecase.MembershipUseCase
}

func NewMembershipController(membershipUseCase usecase.MembershipUseCase) *MembershipController {
	return &MembershipController{
		membershipUseCase: membershipUseCase,
	}
}

func (mc *MembershipController) AddMember(w http.ResponseWriter, r *http.Request) error {
	var req memberdto.AddMemberRequest
	if err := BindRequest(r, &req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errors.WrapValidationError(ve)
		}
		return errors.ErrInvalidInput
	}

	err := mc.membershipUseCase.AddMemberUsecase(r.Context(), *memberdto.ToEntityAddMembersRequest(&req))
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "Member added successfully",
	})
	return nil
}

func (mc *MembershipController) DeleteMember(w http.ResponseWriter, r *http.Request) error {
	var req memberdto.DeleteMemberRequest
	if err := BindRequest(r, &req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return errors.WrapValidationError(ve)
		}
		return errors.ErrInvalidInput
	}

	err := mc.membershipUseCase.DeleteMemberUsecase(r.Context(), *memberdto.ToEntityDeleteMemberRequest(&req), req.CurrentUserID)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Member deleted successfully",
	})
	return nil
}

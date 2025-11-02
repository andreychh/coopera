package web_api

import (
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

// AddMember добавляет пользователя в команду
func (mc *MembershipController) AddMember(w http.ResponseWriter, r *http.Request) {
	var req memberdto.AddMemberRequest

	if err := BindRequest(r, &req); err != nil {
		writeValidationError(w, err)
		return
	}

	err := mc.membershipUseCase.AddMemberUsecase(r.Context(), *memberdto.ToEntityAddMembersRequest(&req))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "Member added successfully",
	})
}

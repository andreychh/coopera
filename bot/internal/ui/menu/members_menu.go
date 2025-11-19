package menu

import (
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
)

func MembersMenu(members []domain.MemberDetails) content.Content {
	return keyboards.Inline(
		content.Text("Team members:"),
		membersMatrix(members).WithRow(buttons.Row(buttons.CallbackButton(
			"Invite member",
			"not_implemented",
		))),
	)
}

func membersMatrix(members []domain.MemberDetails) buttons.ButtonMatrix[buttons.InlineButton] {
	matrix := buttons.Matrix[buttons.InlineButton]()
	for _, member := range members {
		matrix = matrix.WithRow(buttons.Row(memberButton(member)))
	}
	return matrix
}

func memberButton(member domain.MemberDetails) buttons.InlineButton {
	return buttons.CallbackButton(
		member.Name(),
		callbacks.Builder("change_menu").
			With("menu_name", "member").
			With("member_id", strconv.FormatInt(member.ID(), 10)).
			Encode(),
	)
}

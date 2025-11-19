package menu

import (
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
)

func TeamMenu(team domain.TeamDetails) content.Content {
	return keyboards.Inline(
		content.Text(fmt.Sprintf("%s menu:", team.Name())),
		buttons.Matrix(
			buttons.Row(buttons.CallbackButton(
				"Members",
				callbacks.Builder("change_menu").
					With("menu_name", "members").
					With("team_id", strconv.FormatInt(team.ID(), 10)).
					Encode(),
			)),
		),
	)
}

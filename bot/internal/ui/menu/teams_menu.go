package menu

import (
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
)

func TeamsMenu(teams []domain.TeamDetails) content.Content {
	return keyboards.Inline(
		content.Text("Select a team:"),
		teamsMatrix(teams).WithRow(buttons.Row(buttons.CallbackButton(
			"Create team",
			"not_implemented",
		))),
	)
}

func teamsMatrix(teams []domain.TeamDetails) buttons.ButtonMatrix[buttons.InlineButton] {
	matrix := buttons.Matrix[buttons.InlineButton]()
	for _, team := range teams {
		matrix = matrix.WithRow(buttons.Row(teamButton(team)))
	}
	return matrix
}

func teamButton(team domain.TeamDetails) buttons.InlineButton {
	return buttons.CallbackButton(
		team.Name(),
		callbacks.Builder("change_menu").
			With("menu_name", "team").
			With("team_id", strconv.FormatInt(team.ID(), 10)).
			Encode(),
	)
}

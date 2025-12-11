package views

import (
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
)

func MainMenuView() sources.Source[content.Content] {
	return sources.Static[content.Content](
		keyboards.Inline(
			content.Text("Main menu"),
			buttons.Matrix(
				buttons.Row(buttons.CallbackButton("Teams", protocol.ToTeamsMenu())),
				buttons.Row(buttons.CallbackButton("My tasks", "not_implemented")),
			),
		),
	)
}

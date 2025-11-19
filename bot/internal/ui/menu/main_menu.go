package menu

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
)

func MainMenu() content.Content {
	return keyboards.Inline(
		content.Text("Main menu"),
		buttons.Matrix(
			buttons.Row(buttons.CallbackButton("Statistics", "not_implemented")),
			buttons.Row(buttons.CallbackButton(
				"Teams",
				callbacks.Builder("change_menu").With("menu_name", "teams").Encode(),
			)),
		),
	)
}

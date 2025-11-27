package views

import (
	"context"

	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type mainMenuView struct {
}

func (m mainMenuView) Render(_ context.Context, _ telegram.Update) (content.Content, error) {
	return keyboards.Inline(
		content.Text("Main menu"),
		buttons.Matrix(
			buttons.Row(buttons.CallbackButton("Statistics", "not_implemented")),
			buttons.Row(buttons.CallbackButton("Teams", protocol.ToMainMenu())),
		),
	), nil
}

func MainMenu() content.View {
	return mainMenuView{}
}

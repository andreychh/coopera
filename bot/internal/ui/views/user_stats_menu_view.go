package views

import (
	"context"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type userStatsMenuView struct {
	community domain.Community
}

func (u userStatsMenuView) Value(ctx context.Context, update telegram.Update) (content.Content, error) {
	return keyboards.Inline(
		content.Text("User Stats Menu:"),
		buttons.Matrix(
			buttons.Row(buttons.CallbackButton("Main menu", protocol.ToMainMenu())),
		),
	), nil
}

func UserStatsMenu(community domain.Community) sources.Source[content.Content] {
	return userStatsMenuView{
		community: community,
	}
}

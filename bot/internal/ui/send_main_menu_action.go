package ui

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/ui/menu"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type sendMainMenuAction struct {
	bot tg.Bot
}

func (s sendMainMenuAction) Perform(ctx context.Context, update telegram.Update) error {
	id, err := attributes.ChatID(update).Value()
	if err != nil {
		return fmt.Errorf("(%T) getting chat ID: %w", s, err)
	}
	err = s.bot.Chat(id).Send(ctx, menu.MainMenu())
	if err != nil {
		return fmt.Errorf("(%T) sending main menu to chat #%d: %w", s, id, err)
	}
	return nil
}

func SendMainMenu(bot tg.Bot) core.Action {
	return sendMainMenuAction{bot: bot}
}

package domain

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/ui/message"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type sendWelcomeMessageAction struct {
	bot tg.Bot
}

func (s sendWelcomeMessageAction) Perform(ctx context.Context, update telegram.Update) error {
	id, err := attributes.ChatID(update).Value()
	if err != nil {
		return fmt.Errorf("(%T) getting chat ID: %w", s, err)
	}
	err = s.bot.Chat(id).Send(ctx, message.WelcomeMessage())
	if err != nil {
		return err
	}
	return nil
}

func SendWelcomeMessage(bot tg.Bot) core.Action {
	return sendWelcomeMessageAction{bot: bot}
}

package base

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type sendContentAction struct {
	bot     tg.Bot
	content content.Content
}

func (s sendContentAction) Perform(ctx context.Context, update telegram.Update) error {
	id, exists := attrs.ChatID(update).Value()
	if !exists {
		return fmt.Errorf("getting chat ID from update: chat ID not found")
	}
	err := s.bot.Chat(id).Send(ctx, s.content)
	if err != nil {
		return fmt.Errorf("sending content to chat %d: %w", id, err)
	}
	return nil
}

func SendContent(bot tg.Bot, content content.Content) core.Action {
	return sendContentAction{bot: bot, content: content}
}

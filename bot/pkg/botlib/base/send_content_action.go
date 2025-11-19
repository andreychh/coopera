package base

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/base/bot"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type sendContentAction struct {
	bot     bot.Bot
	content content.Content
}

func (s sendContentAction) Perform(ctx context.Context, update telegram.Update) error {
	id, err := attributes.ChatID(update).Value()
	if err != nil {
		return fmt.Errorf("(%T) getting chat ID: %w", s, err)
	}
	err = s.bot.Chat(id).Send(ctx, s.content)
	if err != nil {
		return fmt.Errorf("(%T) sending content to chat id #%d: %w", s, id, err)
	}
	return nil
}

func SendContent(bot bot.Bot, content content.Content) core.Action {
	return sendContentAction{bot: bot, content: content}
}

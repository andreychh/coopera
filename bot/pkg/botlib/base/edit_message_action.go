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

type editMessageAction struct {
	bot  tg.Bot
	view content.View
}

func (s editMessageAction) Perform(ctx context.Context, update telegram.Update) error {
	chatID, exists := attrs.ChatID(update).Value()
	if !exists {
		return fmt.Errorf("getting chat ID from update: chat ID not found")
	}
	messageID, exists := attrs.MessageID(update).Value()
	if !exists {
		return fmt.Errorf("getting message ID from update: message ID not found")
	}
	cnt, err := s.view.Render(ctx, update)
	if err != nil {
		return fmt.Errorf("rendering content for update: %w", err)
	}
	return s.bot.Chat(chatID).Message(messageID).Edit(ctx, cnt)
}

func EditMessage(bot tg.Bot, view content.View) core.Action {
	return editMessageAction{
		bot:  bot,
		view: view,
	}
}

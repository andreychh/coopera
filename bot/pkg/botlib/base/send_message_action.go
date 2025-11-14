package base

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type sendMessageAction struct {
	bot  *telegram.BotAPI
	text string
}

func (s sendMessageAction) Perform(_ context.Context, update telegram.Update) error {
	id, err := attributes.ChatID(update).Value()
	if err != nil {
		return fmt.Errorf("(%T) getting chat ID: %w", s, err)
	}
	message := telegram.NewMessage(id, s.text)
	_, err = s.bot.Send(message)
	if err != nil {
		return fmt.Errorf("(%T->%T) sending message: %w", s, s.bot, err)
	}
	return nil
}

func SendMessage(bot *telegram.BotAPI, text string) core.Action {
	return sendMessageAction{bot: bot, text: text}
}

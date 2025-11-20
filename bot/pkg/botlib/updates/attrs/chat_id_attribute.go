package attrs

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type chatIDAttribute struct {
	update telegram.Update
}

func (c chatIDAttribute) Value() (int64, bool) {
	if c.update.Message != nil {
		return c.update.Message.Chat.ID, true
	}
	if c.update.CallbackQuery != nil && c.update.CallbackQuery.Message != nil {
		return c.update.CallbackQuery.Message.Chat.ID, true
	}
	return 0, false
}

func ChatID(update telegram.Update) Attribute[int64] {
	return chatIDAttribute{update: update}
}

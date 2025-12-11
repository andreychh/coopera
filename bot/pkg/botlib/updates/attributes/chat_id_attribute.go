package attributes

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type chatIDAttribute struct{}

func (c chatIDAttribute) Value(update telegram.Update) (int64, bool) {
	if update.Message != nil {
		return update.Message.Chat.ID, true
	}
	if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
		return update.CallbackQuery.Message.Chat.ID, true
	}
	return 0, false
}

func ChatID() Attribute[int64] {
	return chatIDAttribute{}
}

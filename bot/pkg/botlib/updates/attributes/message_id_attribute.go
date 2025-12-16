package attributes

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type messageIDAttribute struct{}

func (m messageIDAttribute) Value(update telegram.Update) (int, bool) {
	if update.Message != nil {
		return update.Message.MessageID, true
	}
	if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
		return update.CallbackQuery.Message.MessageID, true
	}
	return 0, false
}

func MessageID() Attribute[int] {
	return messageIDAttribute{}
}

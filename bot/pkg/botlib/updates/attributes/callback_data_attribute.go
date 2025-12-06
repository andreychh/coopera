package attributes

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type callbackDataAttribute struct{}

func (c callbackDataAttribute) Value(update telegram.Update) (string, bool) {
	if update.CallbackQuery == nil {
		return "", false
	}
	return update.CallbackQuery.Data, true
}

func CallbackData() Attribute[string] {
	return callbackDataAttribute{}
}

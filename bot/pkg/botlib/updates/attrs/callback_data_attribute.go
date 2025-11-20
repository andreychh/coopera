package attrs

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type callbackDataAttribute struct {
	update telegram.Update
}

func (c callbackDataAttribute) Value() (string, bool) {
	if c.update.CallbackQuery == nil {
		return "", false
	}
	return c.update.CallbackQuery.Data, true
}

func CallbackData(update telegram.Update) Attribute[string] {
	return callbackDataAttribute{update: update}
}

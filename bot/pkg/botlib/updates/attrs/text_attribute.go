package attrs

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type textAttribute struct {
	update telegram.Update
}

func (t textAttribute) Value() (string, bool) {
	if t.update.Message != nil {
		return t.update.Message.Text, true
	}
	return "", false
}

func Text(update telegram.Update) Attribute[string] {
	return textAttribute{update: update}
}

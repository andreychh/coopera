package attributes

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type textAttribute struct{}

func (t textAttribute) Value(update telegram.Update) (string, bool) {
	if update.Message != nil {
		return update.Message.Text, true
	}
	return "", false
}

func Text() Attribute[string] {
	return textAttribute{}
}

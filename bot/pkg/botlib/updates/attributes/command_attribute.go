package attributes

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type commandAttribute struct{}

func (c commandAttribute) Value(update telegram.Update) (string, bool) {
	if update.Message == nil || !update.Message.IsCommand() {
		return "", false
	}
	return update.Message.Command(), true
}

func Command() Attribute[string] {
	return commandAttribute{}
}

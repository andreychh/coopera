package attrs

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type commandAttribute struct {
	update telegram.Update
}

func (c commandAttribute) Value() (string, bool) {
	if c.update.Message == nil {
		return "", false
	}
	if !c.update.Message.IsCommand() {
		return "", false
	}
	return c.update.Message.Command(), true
}

func Command(update telegram.Update) attrs.Attribute[string] {
	return commandAttribute{update: update}
}

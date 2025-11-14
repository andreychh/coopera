package attributes

import (
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type commandAttribute struct {
	update telegram.Update
}

func (c commandAttribute) Exists() bool {
	return c.update.Message != nil && c.update.Message.IsCommand()
}

func (c commandAttribute) Value() (string, error) {
	if c.update.Message != nil && c.update.Message.IsCommand() {
		return c.update.Message.Command(), nil
	}
	return "", fmt.Errorf("(%T) %w", c, ErrAttributeNotFound)
}

func Command(update telegram.Update) Attribute[string] {
	return commandAttribute{update: update}
}

package attributes

import (
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type textAttribute struct {
	update telegram.Update
}

func (c textAttribute) Exists() bool {
	return c.update.Message != nil
}

func (c textAttribute) Value() (string, error) {
	if c.update.Message != nil {
		return c.update.Message.Text, nil
	}
	return "", fmt.Errorf("(%T) %w", c, ErrAttributeNotFound)
}

func Text(update telegram.Update) Attribute[string] {
	return textAttribute{update: update}
}

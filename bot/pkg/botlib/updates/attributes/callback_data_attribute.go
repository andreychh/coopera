package attributes

import (
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type callbackDataAttribute struct {
	update telegram.Update
}

func (c callbackDataAttribute) Exists() bool {
	return c.update.CallbackQuery != nil
}

func (c callbackDataAttribute) Value() (string, error) {
	if c.update.CallbackQuery != nil {
		return c.update.CallbackQuery.Data, nil
	}
	return "", fmt.Errorf("(%T) %w", c, ErrAttributeNotFound)
}

func CallbackData(update telegram.Update) Attribute[string] {
	return callbackDataAttribute{update: update}
}

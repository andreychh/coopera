package attributes

import (
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type chatIDAttribute struct {
	update telegram.Update
}

func (c chatIDAttribute) Exists() bool {
	return (c.update.Message != nil) ||
		(c.update.CallbackQuery != nil && c.update.CallbackQuery.Message != nil)
}

func (c chatIDAttribute) Value() (int64, error) {
	if c.update.Message != nil {
		return c.update.Message.Chat.ID, nil
	}
	if c.update.CallbackQuery != nil && c.update.CallbackQuery.Message != nil {
		return c.update.CallbackQuery.Message.Chat.ID, nil
	}
	return 0, fmt.Errorf("(%T) %w", c, ErrAttributeNotFound)
}

func ChatID(update telegram.Update) Attribute[int64] {
	return chatIDAttribute{update: update}
}

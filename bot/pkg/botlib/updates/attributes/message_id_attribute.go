package attributes

import telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type messageIDAttribute struct {
	update telegram.Update
}

func (m messageIDAttribute) Exists() bool {
	return (m.update.Message != nil) ||
		(m.update.CallbackQuery != nil && m.update.CallbackQuery.Message != nil)
}

func (m messageIDAttribute) Value() (int, error) {
	if m.update.Message != nil {
		return m.update.Message.MessageID, nil
	}

	if m.update.CallbackQuery != nil && m.update.CallbackQuery.Message != nil {
		return m.update.CallbackQuery.Message.MessageID, nil
	}
	return 0, ErrAttributeNotFound
}

func MessageID(update telegram.Update) Attribute[int] {
	return messageIDAttribute{update: update}
}

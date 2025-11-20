package attrs

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type messageIDAttribute struct {
	update telegram.Update
}

func (m messageIDAttribute) Value() (int, bool) {
	if m.update.Message != nil {
		return m.update.Message.MessageID, true
	}
	if m.update.CallbackQuery != nil && m.update.CallbackQuery.Message != nil {
		return m.update.CallbackQuery.Message.MessageID, true
	}
	return 0, false
}

func MessageID(update telegram.Update) attrs.Attribute[int] {
	return messageIDAttribute{update: update}
}

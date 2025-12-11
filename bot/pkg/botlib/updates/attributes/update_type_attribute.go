package attributes

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type updateTypeAttribute struct{}

func (u updateTypeAttribute) Value(update telegram.Update) (updates.UpdateType, bool) {
	if update.Message != nil {
		return updates.UpdateTypeMessage, true
	}
	if update.CallbackQuery != nil {
		return updates.UpdateTypeCallbackQuery, true
	}
	return updates.UpdateTypeUnknown, false
}

func UpdateType() Attribute[updates.UpdateType] {
	return updateTypeAttribute{}
}

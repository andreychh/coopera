package attrs

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type updateTypeAttribute struct {
	update telegram.Update
}

func (u updateTypeAttribute) Value() (updates.UpdateType, bool) {
	if u.update.Message != nil {
		return updates.UpdateTypeMessage, true
	}
	if u.update.CallbackQuery != nil {
		return updates.UpdateTypeCallbackQuery, true
	}
	return updates.UpdateTypeUnknown, false
}

func UpdateType(update telegram.Update) Attribute[updates.UpdateType] {
	return updateTypeAttribute{update: update}
}

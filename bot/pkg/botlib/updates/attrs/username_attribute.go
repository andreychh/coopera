package attrs

import telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type usernameAttribute struct {
	update telegram.Update
}

func (u usernameAttribute) Value() (value string, exists bool) {
	if u.update.Message != nil && u.update.Message.From != nil {
		return u.update.Message.From.UserName, true
	}
	if u.update.EditedMessage != nil && u.update.EditedMessage.From != nil {
		return u.update.EditedMessage.From.UserName, true
	}
	if u.update.CallbackQuery != nil && u.update.CallbackQuery.From != nil {
		return u.update.CallbackQuery.From.UserName, true
	}
	return "", false
}

func Username(update telegram.Update) Attribute[string] {
	return usernameAttribute{update: update}
}

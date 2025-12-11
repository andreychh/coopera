package attributes

import telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type usernameAttribute struct{}

func (u usernameAttribute) Value(update telegram.Update) (string, bool) {
	if update.Message != nil && update.Message.From != nil {
		return update.Message.From.UserName, true
	}
	if update.EditedMessage != nil && update.EditedMessage.From != nil {
		return update.EditedMessage.From.UserName, true
	}
	if update.CallbackQuery != nil && update.CallbackQuery.From != nil {
		return update.CallbackQuery.From.UserName, true
	}
	return "", false
}

func Username() Attribute[string] {
	return usernameAttribute{}
}

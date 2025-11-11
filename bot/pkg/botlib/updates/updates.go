package updates

import (
	"errors"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ErrNoChatID = errors.New("no chat ID available in update")
var ErrNoText = errors.New("no text available in update")
var ErrNoCommand = errors.New("no command available in update")

func ChatID(update telegram.Update) (id int64, available bool) {
	if update.Message != nil {
		return update.Message.Chat.ID, true
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID, true
	}
	return 0, false
}

func Text(update telegram.Update) (text string, available bool) {
	if update.Message != nil {
		return update.Message.Text, true
	}
	return "", false
}

func Command(update telegram.Update) (command string, available bool) {
	if update.Message == nil {
		return "", false
	}
	if !update.Message.IsCommand() {
		return "", false
	}
	return update.Message.Command(), true
}

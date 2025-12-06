package attributes

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type chatTypeAttribute struct{}

func (c chatTypeAttribute) Value(update telegram.Update) (updates.ChatType, bool) {
	if update.Message != nil && update.Message.Chat != nil {
		return updates.ChatType(update.Message.Chat.Type), true
	}
	if update.EditedMessage != nil && update.EditedMessage.Chat != nil {
		return updates.ChatType(update.EditedMessage.Chat.Type), true
	}
	if update.ChannelPost != nil && update.ChannelPost.Chat != nil {
		return updates.ChatType(update.ChannelPost.Chat.Type), true
	}
	if update.EditedChannelPost != nil && update.EditedChannelPost.Chat != nil {
		return updates.ChatType(update.EditedChannelPost.Chat.Type), true
	}
	if update.CallbackQuery != nil && update.CallbackQuery.Message != nil && update.CallbackQuery.Message.Chat != nil {
		return updates.ChatType(update.CallbackQuery.Message.Chat.Type), true
	}
	return updates.ChatTypeUnknown, false
}

func ChatType() Attribute[updates.ChatType] {
	return chatTypeAttribute{}
}

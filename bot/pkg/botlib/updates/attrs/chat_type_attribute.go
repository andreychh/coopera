package attrs

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type chatTypeAttribute struct {
	telegram.Update
}

func (c chatTypeAttribute) Value() (value updates.ChatType, exists bool) {
	if c.Message != nil && c.Message.Chat != nil {
		return updates.ChatType(c.Message.Chat.Type), true
	}
	if c.EditedMessage != nil && c.EditedMessage.Chat != nil {
		return updates.ChatType(c.EditedMessage.Chat.Type), true
	}
	if c.ChannelPost != nil && c.ChannelPost.Chat != nil {
		return updates.ChatType(c.ChannelPost.Chat.Type), true
	}
	if c.EditedChannelPost != nil && c.EditedChannelPost.Chat != nil {
		return updates.ChatType(c.EditedChannelPost.Chat.Type), true
	}
	if c.CallbackQuery != nil && c.CallbackQuery.Message != nil && c.CallbackQuery.Message.Chat != nil {
		return updates.ChatType(c.CallbackQuery.Message.Chat.Type), true
	}
	return updates.ChatTypeUnknown, false
}

func ChatType(update telegram.Update) Attribute[updates.ChatType] {
	return chatTypeAttribute{Update: update}
}

package tg

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
)

func ChatWithID(bot Bot, id sources.Source[int64]) sources.Source[Chat] {
	return sources.PureMap(id,
		func(id int64) Chat {
			return bot.Chat(id)
		},
	)
}

func CurrentChat(bot Bot) sources.Source[Chat] {
	return ChatWithID(bot, sources.Required(attributes.ChatID()))
}

func MessageWithID(chat sources.Source[Chat], id sources.Source[int]) sources.Source[Message] {
	return sources.PureZip(chat, id,
		func(chat Chat, id int) Message {
			return chat.Message(id)
		},
	)
}

func CurrentMessage(bot Bot) sources.Source[Message] {
	return MessageWithID(CurrentChat(bot), sources.Required(attributes.MessageID()))
}

package bot

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/base/chat"
	"github.com/andreychh/coopera-bot/pkg/botlib/base/client"
)

type Bot interface {
	Chat(id int64) chat.Chat
}
type bot struct {
	dataSource client.Client
}

func (b bot) Chat(id int64) chat.Chat {
	return chat.New(id, b.dataSource)
}

func New(dataSource client.Client) Bot {
	return bot{
		dataSource: dataSource,
	}
}

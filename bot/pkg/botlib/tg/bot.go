package tg

import "github.com/andreychh/coopera-bot/pkg/botlib/transport"

type bot struct {
	dataSource transport.Client
}

func (b bot) Chat(id int64) Chat {
	return NewChat(id, b.dataSource)
}

func NewBot(dataSource transport.Client) Bot {
	return bot{
		dataSource: dataSource,
	}
}

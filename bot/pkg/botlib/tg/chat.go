package tg

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/transport"
)

type chat struct {
	id         int64
	dataSource transport.Client
}

func (c chat) Send(ctx context.Context, cnt content.Content) error {
	cnt = content.WithChatID(cnt, c.id)
	payload, err := cnt.Structure().Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal content structure: %w", err)
	}
	_, err = c.dataSource.Execute(ctx, cnt.Method(), payload)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	return nil
}

func (c chat) Message(id int) Message {
	return NewMessage(c.id, id, c.dataSource)
}

func NewChat(id int64, dataSource transport.Client) Chat {
	return chat{
		id:         id,
		dataSource: dataSource,
	}
}

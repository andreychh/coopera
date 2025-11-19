package chat

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/base/client"
	"github.com/andreychh/coopera-bot/pkg/botlib/base/message"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
)

type Chat interface {
	Send(ctx context.Context, cnt content.Content) error
	Message(id int64) message.Message
}

type chat struct {
	id         int64
	dataSource client.Client
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

func (c chat) Message(id int64) message.Message {
	return message.New(c.id, id, c.dataSource)
}

func New(id int64, dataSource client.Client) Chat {
	return chat{
		id:         id,
		dataSource: dataSource,
	}
}

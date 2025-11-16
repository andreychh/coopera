package chat

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/base/client"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
)

type Chat interface {
	Send(ctx context.Context, cnt content.Content) error
}

type chat struct {
	id         int64
	dataSource client.Client
}

func (c chat) Send(ctx context.Context, cnt content.Content) error {
	cnt = content.To(cnt, c.id)
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

func New(id int64, dataSource client.Client) Chat {
	return chat{
		id:         id,
		dataSource: dataSource,
	}
}

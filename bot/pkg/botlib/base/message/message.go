package message

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/base/client"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type Message interface {
	Edit(ctx context.Context, content content.Content) error
	Delete(ctx context.Context) error
}

type message struct {
	chatID     int64
	messageID  int64
	dataSource client.Client
}

func (m message) Edit(ctx context.Context, cnt content.Content) error {
	cnt = content.WithMessageID(content.WithChatID(cnt, m.chatID), m.messageID)
	payload, err := cnt.Structure().Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal content structure: %w", err)
	}
	_, err = m.dataSource.Execute(ctx, "editMessageText", payload)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	return nil
}

func (m message) Delete(ctx context.Context) error {
	payload, err := json.Object(json.Fields{
		"chat_id":    json.I64(m.chatID),
		"message_id": json.I64(m.messageID),
	}).Marshal()
	if err != nil {
		return fmt.Errorf("marshaling delete message payload: %w", err)
	}
	_, err = m.dataSource.Execute(ctx, "deleteMessage", payload)
	if err != nil {
		return fmt.Errorf("deleting message: %w", err)
	}
	return nil
}

func New(chatID int64, messageID int64, dataSource client.Client) Message {
	return message{
		chatID:     chatID,
		messageID:  messageID,
		dataSource: dataSource,
	}
}

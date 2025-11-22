package tg

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg/transport"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

var ErrMessageCannotBeEdited = fmt.Errorf("message cannot be edited")

type message struct {
	chatID     int64
	messageID  int
	dataSource transport.Client
}

// TODO: "editMessageText" is suitable only for text messages. Need to handle other types of messages.
// TODO: Handle ErrMessageCannotBeEdited when editing is not possible.
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
		"message_id": json.Int(m.messageID),
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

func NewMessage(chatID int64, messageID int, dataSource transport.Client) Message {
	return message{
		chatID:     chatID,
		messageID:  messageID,
		dataSource: dataSource,
	}
}

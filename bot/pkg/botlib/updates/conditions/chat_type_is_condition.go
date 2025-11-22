package conditions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type chatTypeIsCondition struct {
	chatType updates.ChatType
}

func (c chatTypeIsCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	chatType, exists := attrs.ChatType(update).Value()
	if !exists {
		return false, fmt.Errorf("getting chat type: chat type not found")
	}
	return chatType == c.chatType, nil
}

func ChatTypeIs(chatType updates.ChatType) core.Condition {
	return chatTypeIsCondition{chatType: chatType}
}

package conditions

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type chatTypeIsCondition struct {
	target updates.ChatType
}

func (c chatTypeIsCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	chatType, exists := attributes.ChatType().Value(update)
	if !exists {
		return false, nil
	}
	return chatType == c.target, nil
}

func ChatTypeIs(target updates.ChatType) core.Condition {
	return chatTypeIsCondition{target: target}
}

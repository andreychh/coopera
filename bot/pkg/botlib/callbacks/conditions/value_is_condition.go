package conditions

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type valueIsCondition struct {
	key   string
	value string
}

func (p valueIsCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	callbackData, exists := attributes.CallbackData().Value(update)
	if !exists {
		return false, nil
	}
	value, exists := callbacks.IncomingData(callbackData).Value(p.key)
	if !exists {
		return false, nil
	}
	return value == p.value, nil
}

func ValueIs(key string, value string) core.Condition {
	return valueIsCondition{key: key, value: value}
}

package conditions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type valueIsCondition struct {
	key   string
	value string
}

func (p valueIsCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	callbackData, exists := attrs.CallbackData(update).Value()
	if !exists {
		return false, fmt.Errorf("getting callback data from update: callback data not found")
	}
	value, exists := callbacks.IncomingData(callbackData).Value(p.key)
	if !exists {
		return false, fmt.Errorf("getting parameter %q from callback data: parameter not found", p.key)
	}
	return value == p.value, nil
}

// ValueIs requires attrs.CallbackData
func ValueIs(key string, value string) core.Condition {
	return valueIsCondition{key: key, value: value}
}

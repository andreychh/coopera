package conditions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type prefixIsCondition struct {
	prefix string
}

func (p prefixIsCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	callbackData, exists := attrs.CallbackData(update).Value()
	if !exists {
		return false, fmt.Errorf("getting callback data from update: callback data not found")
	}
	return callbacks.IncomingData(callbackData).Prefix() == p.prefix, nil
}

// PrefixIs requires attrs.CallbackData
func PrefixIs(prefix string) core.Condition {
	return prefixIsCondition{prefix: prefix}
}

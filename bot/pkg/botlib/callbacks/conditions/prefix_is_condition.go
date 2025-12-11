package conditions

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type prefixIsCondition struct {
	prefix string
}

func (p prefixIsCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	callbackData, exists := attributes.CallbackData().Value(update)
	if !exists {
		return false, nil
	}
	return callbacks.IncomingData(callbackData).Prefix() == p.prefix, nil
}

func PrefixIs(prefix string) core.Condition {
	return prefixIsCondition{prefix: prefix}
}

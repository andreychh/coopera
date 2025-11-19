package callbacks

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type prefixIsCondition struct {
	prefix string
}

func (p prefixIsCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	callbackData, err := attributes.CallbackData(update).Value()
	if err != nil {
		return false, fmt.Errorf("(%T) getting callback data: %w", p, err)
	}
	prefix, err := PrefixedData(callbackData).Prefix()
	if err != nil {
		return false, fmt.Errorf("(%T) getting prefix from callback data: %w", p, err)
	}
	return prefix == p.prefix, nil
}

func PrefixIs(s string) core.Condition {
	return prefixIsCondition{prefix: s}
}

func SafePrefixIs(s string) core.Condition {
	return composition.All(updates.HasCallbackData(), PrefixIs(s))
}

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

type paramIsCondition struct {
	key   string
	value string
}

func (p paramIsCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	callbackData, err := attributes.CallbackData(update).Value()
	if err != nil {
		return false, fmt.Errorf("(%T) getting callback data: %w", p, err)
	}
	value, err := PrefixedData(callbackData).Value(p.key)
	if err != nil {
		return false, fmt.Errorf("(%T) getting param %q from callback data: %w", p, p.key, err)
	}
	return value == p.value, nil
}

func ParamIs(key, value string) core.Condition {
	return paramIsCondition{key: key, value: value}
}

func SafeParamIs(key, value string) core.Condition {
	return composition.All(updates.HasCallbackData(), ParamIs(key, value))
}

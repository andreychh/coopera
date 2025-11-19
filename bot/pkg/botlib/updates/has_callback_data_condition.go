package updates

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type hasCallbackDataCondition struct{}

func (h hasCallbackDataCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	return attributes.CallbackData(update).Exists(), nil
}

func HasCallbackData() core.Condition {
	return hasCallbackDataCondition{}
}

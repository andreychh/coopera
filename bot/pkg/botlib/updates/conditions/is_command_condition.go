package conditions

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type isCommandCondition struct{}

func (i isCommandCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	_, exists := attributes.Command().Value(update)
	return exists, nil
}

func IsCommand() core.Condition {
	return isCommandCondition{}
}

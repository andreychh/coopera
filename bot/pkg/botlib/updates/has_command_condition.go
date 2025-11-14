package updates

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type hasCommandCondition struct{}

func (c hasCommandCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	return attributes.Command(update).Exists(), nil
}

func HasCommand() core.Condition {
	return hasCommandCondition{}
}

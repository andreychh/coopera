package updates

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type hasCommandCondition struct{}

func (c hasCommandCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	_, available := Command(update)
	return available, nil
}

func HasCommand() core.Condition {
	return hasCommandCondition{}
}

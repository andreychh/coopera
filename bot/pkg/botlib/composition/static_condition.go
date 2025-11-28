package composition

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type staticCondition struct {
	value bool
}

func (a staticCondition) Holds(_ context.Context, _ telegram.Update) (bool, error) {
	return a.value, nil
}

func Always() core.Condition {
	return staticCondition{
		value: true,
	}
}

func Never() core.Condition {
	return staticCondition{
		value: false,
	}
}

package updates

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type hasTextCondition struct{}

func (c hasTextCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	_, available := Text(update)
	return available, nil
}

func HasText() core.Condition {
	return hasTextCondition{}
}

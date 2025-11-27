package composition

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type nothingAction struct{}

func (n nothingAction) Perform(_ context.Context, _ telegram.Update) error {
	return nil
}

func Nothing() core.Action {
	return nothingAction{}
}

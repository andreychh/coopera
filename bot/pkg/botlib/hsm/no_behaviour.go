package hsm

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type noBehavior struct{}

func (n noBehavior) Prompt(_ context.Context, _ telegram.Update) error {
	return nil
}

func (n noBehavior) React(_ context.Context, _ telegram.Update) (Decision, error) {
	return Pass(), nil
}

func (n noBehavior) Cleanup(_ context.Context, _ telegram.Update) error {
	return nil
}

func NoBehavior() Behavior {
	return noBehavior{}
}

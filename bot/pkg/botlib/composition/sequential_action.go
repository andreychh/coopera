package composition

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type sequentialAction struct {
	actions []core.Action
}

func (s sequentialAction) Perform(ctx context.Context, update telegram.Update) error {
	for i, act := range s.actions {
		err := act.Perform(ctx, update)
		if err != nil {
			return fmt.Errorf("(%T->%T) performing action #%d: %w", s, act, i, err)
		}
	}
	return nil
}

func Sequential(actions ...core.Action) core.Action {
	return sequentialAction{actions: actions}
}

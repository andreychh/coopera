package composition

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type anyCondition struct {
	conditions []core.Condition
}

func (a anyCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	for i, cnd := range a.conditions {
		held, err := cnd.Holds(ctx, update)
		if err != nil {
			return false, fmt.Errorf("(%T->%T) checking condition #%d: %w", a, cnd, i, err)
		}
		if held {
			return true, nil
		}
	}
	return false, nil
}

func Any(conditions ...core.Condition) core.Condition {
	return anyCondition{conditions: conditions}
}

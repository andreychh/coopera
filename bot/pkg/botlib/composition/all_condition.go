package composition

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type allCondition struct {
	conditions []core.Condition
}

func (a allCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	for i, cnd := range a.conditions {
		held, err := cnd.Holds(ctx, update)
		if err != nil {
			return false, fmt.Errorf("(%T->%T) checking condition #%d: %w", a, cnd, i, err)
		}
		if !held {
			return false, nil
		}
	}
	return true, nil
}

func All(conditions ...core.Condition) core.Condition {
	return allCondition{conditions: conditions}
}

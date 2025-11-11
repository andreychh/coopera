package composition

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type notCondition struct {
	origin core.Condition
}

func (n notCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	held, err := n.origin.Holds(ctx, update)
	if err != nil {
		return false, fmt.Errorf("(%T->%T) checking condition: %w", n, n.origin, err)
	}
	return !held, nil
}

func Not(origin core.Condition) core.Condition {
	return notCondition{origin: origin}
}

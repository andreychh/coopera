package routing

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type conditionalClause struct {
	origin    core.Clause
	condition core.Condition
}

func (c conditionalClause) TryExecute(ctx context.Context, update telegram.Update) (bool, error) {
	held, err := c.condition.Holds(ctx, update)
	if err != nil {
		return false, fmt.Errorf("(%T->%T) checking condition: %w", c, c.condition, err)
	}
	if !held {
		return false, nil
	}
	executed, err := c.origin.TryExecute(ctx, update)
	if err != nil {
		return false, fmt.Errorf("(%T->%T) trying to execute clause: %w", c, c.origin, err)
	}
	return executed, nil
}

func If(condition core.Condition, origin core.Clause) core.Clause {
	return conditionalClause{
		origin:    origin,
		condition: condition,
	}
}

func TerminalIf(condition core.Condition, action core.Action) core.Clause {
	return If(condition, Terminal(action))
}

package routing

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type terminalClause struct {
	action core.Action
}

func (t terminalClause) TryExecute(ctx context.Context, update telegram.Update) (bool, error) {
	err := t.action.Perform(ctx, update)
	if err != nil {
		return false, fmt.Errorf("(%T->%T) performing action: %w", t, t.action, err)
	}
	return true, nil
}

func Terminal(action core.Action) core.Clause {
	return terminalClause{action: action}
}

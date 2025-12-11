package composition

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type clauseAction struct {
	clause   core.Clause
	fallback core.Action
}

func (c clauseAction) Perform(ctx context.Context, u telegram.Update) error {
	handled, err := c.clause.TryExecute(ctx, u)
	if err != nil {
		return fmt.Errorf("logic execution failed: %w", err)
	}
	if !handled {
		return c.fallback.Perform(ctx, u)
	}
	return nil
}

func FallbackRun(c core.Clause, fallback core.Action) core.Action {
	return clauseAction{
		clause:   c,
		fallback: fallback,
	}
}

func Run(c core.Clause) core.Action {
	return clauseAction{
		clause:   c,
		fallback: Nothing(),
	}
}

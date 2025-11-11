package routing

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type firstExecutedClause struct {
	clauses []core.Clause
}

func (f firstExecutedClause) TryExecute(ctx context.Context, update telegram.Update) (bool, error) {
	for i, cls := range f.clauses {
		executed, err := cls.TryExecute(ctx, update)
		if err != nil {
			return false, fmt.Errorf("(%T, %T): executing clause #%d: %w", f, cls, i, err)
		}
		if executed {
			return true, nil
		}
	}
	return false, nil
}

func FirstExecuted(clauses ...core.Clause) core.Clause {
	return firstExecutedClause{clauses: clauses}
}

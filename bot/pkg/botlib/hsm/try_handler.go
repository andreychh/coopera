package hsm

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/routing"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type tryHandler struct {
	clause   core.Clause
	decision Decision
}

func (j tryHandler) Handle(ctx context.Context, u telegram.Update) (Decision, error) {
	executed, err := j.clause.TryExecute(ctx, u)
	if err != nil {
		return nil, err
	}
	if executed {
		return j.decision, nil
	}
	return Pass(), nil
}

func Try(clause core.Clause, decision Decision) Handler {
	return tryHandler{
		clause:   clause,
		decision: decision,
	}
}

func TryAction(condition core.Condition, action core.Action, decision Decision) Handler {
	return Try(routing.TerminalIf(condition, action), decision)
}

package base

import (
	"context"
	"errors"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var errRecoveredPanic = errors.New("recovered from panic during clause execution")

type recoveryClause struct {
	origin core.Clause
}

func (a recoveryClause) TryExecute(ctx context.Context, update telegram.Update) (executed bool, err error) {
	defer func() {
		value := recover()
		if value != nil {
			panicErr := fmt.Errorf("%w, value is %s", errRecoveredPanic, value)
			if err != nil {
				err = errors.Join(err, panicErr)
			} else {
				err = panicErr
			}
			executed = false
		}
	}()
	executed, err = a.origin.TryExecute(ctx, update)
	if err != nil {
		return false, fmt.Errorf("(%T->%T) executing clause: %w", a, a.origin, err)
	}
	return executed, nil
}

func Recover(origin core.Clause) core.Clause {
	return recoveryClause{origin: origin}
}

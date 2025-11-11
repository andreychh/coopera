package logging

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type loggingClause struct {
	origin core.Clause
	logger *slog.Logger
}

func (l loggingClause) TryExecute(ctx context.Context, update telegram.Update) (bool, error) {
	executed, err := l.origin.TryExecute(ctx, update)
	if err != nil {
		l.logger.ErrorContext(ctx, "clause execution error",
			"clause_type", fmt.Sprintf("%T", l.origin),
			"error", err,
		)
		return false, err
	}
	if executed {
		l.logger.InfoContext(ctx, "clause executed",
			"clause_type", fmt.Sprintf("%T", l.origin),
		)
	} else {
		l.logger.InfoContext(ctx, "clause not executed",
			"clause_type", fmt.Sprintf("%T", l.origin),
		)
	}
	return executed, nil
}

func LoggingClause(origin core.Clause, logger *slog.Logger) core.Clause {
	return loggingClause{origin: origin, logger: logger}
}

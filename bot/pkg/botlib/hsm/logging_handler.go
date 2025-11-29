package hsm

import (
	"context"
	"log/slog"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type loggingHandler struct {
	logger *slog.Logger
	origin Handler
}

func (l loggingHandler) Handle(ctx context.Context, u telegram.Update) (Decision, error) {
	decision, err := l.origin.Handle(ctx, u)
	if err != nil {
		l.logger.Error("handler execution failed",
			"error", err,
			"update_id", u.UpdateID,
		)
		return nil, err
	}
	if decision.Handled() {
		l.logger.Info("handler processed update",
			"handled", true,
			"target_id", decision.Next(),
			"update_id", u.UpdateID,
		)
	} else {
		l.logger.Debug("handler ignored update",
			"handled", false,
			"update_id", u.UpdateID,
		)
	}
	return decision, nil
}

func Logging(logger *slog.Logger, origin Handler) Handler {
	return loggingHandler{
		logger: logger,
		origin: origin,
	}
}

func DefaultLogging(origin Handler) Handler {
	return Logging(slog.Default(), origin)
}

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package actions

import (
	"context"
	"log/slog"

	"github.com/andreychh/coopera/pkg/bot/api"
)

type loggingAction struct {
	logger *slog.Logger
}

func (a loggingAction) Execute(ctx context.Context, update api.Update) error {
	a.logger.InfoContext(ctx, "update received",
		slog.Int("update.id", int(update.UpdateID)),
	)
	return nil
}

func LoggingAction(logger *slog.Logger) Action {
	return loggingAction{
		logger: logger,
	}
}

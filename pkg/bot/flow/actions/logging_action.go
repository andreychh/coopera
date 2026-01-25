// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package actions

import (
	"context"
	"log/slog"

	"github.com/andreychh/coopera/pkg/bot/api"
)

type LoggingAction struct {
	logger *slog.Logger
}

func NewLoggingAction(logger *slog.Logger) LoggingAction {
	return LoggingAction{
		logger: logger,
	}
}

func (a LoggingAction) Execute(ctx context.Context, update api.Update) error {
	a.logger.InfoContext(ctx, "update received",
		slog.Int("update.id", int(update.UpdateID)),
	)
	return nil
}

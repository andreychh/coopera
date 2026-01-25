// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package transport

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/andreychh/coopera/pkg/bot/api"
)

type LoggingClient struct {
	origin TelegramClient
	logger *slog.Logger
}

func NewLoggingClient(origin TelegramClient, logger *slog.Logger) LoggingClient {
	return LoggingClient{
		origin: origin,
		logger: logger,
	}
}

func (c LoggingClient) SendRequest(ctx context.Context, method api.Method, requestBody, responseBody any) error {
	start := time.Now()
	err := c.origin.SendRequest(ctx, method, requestBody, responseBody)
	duration := time.Since(start)
	var apiErr *api.Error
	if errors.As(err, &apiErr) {
		c.logger.ErrorContext(ctx, "request failed",
			slog.String("method", string(method)),
			slog.Duration("duration", duration),
			slog.String("error", err.Error()),
			slog.String("request", fmt.Sprintf("%+v", requestBody)),
			slog.Any("envelope", api.LogEnvelope(apiErr.Envelope())),
		)
		return err
	}
	if err != nil {
		c.logger.ErrorContext(ctx, "request failed",
			slog.String("method", string(method)),
			slog.Duration("duration", duration),
			slog.String("error", err.Error()),
			slog.String("request", fmt.Sprintf("%+v", requestBody)),
		)
		return err
	}
	c.logger.InfoContext(ctx, "request succeeded",
		slog.String("method", string(method)),
		slog.Duration("duration", duration),
		slog.String("response.type", fmt.Sprintf("%T", responseBody)),
	)
	return nil
}

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/andreychh/coopera/pkg/bot/app"
	"github.com/andreychh/coopera/pkg/bot/endpoints"
	"github.com/andreychh/coopera/pkg/bot/flow/actions"
	"github.com/andreychh/coopera/pkg/bot/flow/updates"
	"github.com/andreychh/coopera/pkg/bot/transport"
)

func main() {
	token, exists := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if !exists {
		panic("TELEGRAM_BOT_TOKEN environment variable is not set")
	}
	logger := Logger()
	client := transport.NewLoggingClient(TelegramClient(token), logger)
	action := actions.LoggingAction(logger)
	application := app.SingleWorkerApp(updates.LongPollingSource(endpoints.GetUpdates(client)), action)
	err := application.Run(context.Background())
	if err != nil {
		panic(err)
	}
}

func TelegramClient(token string) transport.TelegramClient {
	return transport.NewClient(token, HTTPClient())
}

func HTTPClient() *http.Client {
	return http.DefaultClient
}

func Logger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

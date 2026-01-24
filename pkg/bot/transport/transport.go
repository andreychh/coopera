// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package transport

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/pkg/bot/api"
)

type TelegramClient interface {
	SendRequest(ctx context.Context, method api.Method, requestBody, responseBody any) error
}

type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

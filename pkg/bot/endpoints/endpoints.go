// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package endpoints

import (
	"context"

	"github.com/andreychh/coopera/pkg/bot/api"
)

type TelegramClient interface {
	SendRequest(ctx context.Context, method api.Method, reqBody, respBody any) error
}

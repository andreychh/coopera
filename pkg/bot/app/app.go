// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package app

import (
	"context"

	"github.com/andreychh/coopera/pkg/bot/api"
)

type UpdateSource interface {
	Updates(ctx context.Context) <-chan api.Update
}

type Action interface {
	Execute(ctx context.Context, update api.Update) error
}

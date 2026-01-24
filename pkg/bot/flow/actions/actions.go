// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package actions

import (
	"context"

	"github.com/andreychh/coopera/pkg/bot/api"
)

type Action interface {
	Execute(ctx context.Context, update api.Update) error
}

type PartialAction interface {
	TryExecute(ctx context.Context, update api.Update) (executed bool, err error)
}

type Predicate interface {
	Evaluate(ctx context.Context, update api.Update) (result bool, err error)
}

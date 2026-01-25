// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package app

import (
	"context"
)

type SingleWorkerApp struct {
	source UpdateSource
	action Action
}

func NewSingleWorkerApp(source UpdateSource, action Action) SingleWorkerApp {
	return SingleWorkerApp{
		source: source,
		action: action,
	}
}

func (a SingleWorkerApp) Run(ctx context.Context) error {
	for update := range a.source.Updates(ctx) {
		_ = a.action.Execute(ctx, update)
	}
	return ctx.Err()
}

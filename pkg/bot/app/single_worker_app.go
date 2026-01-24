// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package app

import (
	"context"

	"github.com/andreychh/coopera/pkg/bot/flow/actions"
	"github.com/andreychh/coopera/pkg/bot/flow/updates"
)

type singleWorkerApp struct {
	source updates.UpdateSource
	action actions.Action
}

func (a singleWorkerApp) Run(ctx context.Context) error {
	for update := range a.source.Updates(ctx) {
		_ = a.action.Execute(ctx, update)
	}
	return ctx.Err()
}

func SingleWorkerApp(source updates.UpdateSource, action actions.Action) App {
	return singleWorkerApp{
		source: source,
		action: action,
	}
}

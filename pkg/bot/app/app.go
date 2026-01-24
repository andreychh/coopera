// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package app

import "context"

type App interface {
	Run(ctx context.Context) error
}

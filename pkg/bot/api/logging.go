// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"log/slog"

	"github.com/andreychh/coopera/pkg/ptr"
)

func LogEnvelope(e Envelope) slog.Value {
	return slog.GroupValue(
		slog.Bool("ok", e.Ok),
		slog.Int("error_code", int(ptr.ValueOrDefault(e.ErrorCode, -1))),
		slog.String("description", ptr.ValueOrDefault(e.Description, "unknown")),
	)
}

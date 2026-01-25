// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package updates

import (
	"context"
)

type Endpoint[Req, Resp any] interface {
	Call(ctx context.Context, req Req) (Resp, error)
}

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"errors"
)

func AsError(err error) *Error {
	var e *Error
	if errors.As(err, e) {
		return e
	}
	return nil
}

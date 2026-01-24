// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package transport

import (
	"strings"
)

type mappedError struct {
	origin   error
	old, new string
}

func (e mappedError) Error() string {
	return strings.ReplaceAll(e.origin.Error(), e.old, e.new)
}

func (e mappedError) Unwrap() error {
	return e.origin
}

func MappedError(origin error, old, new string) error {
	return &mappedError{
		origin: origin,
		old:    old,
		new:    new,
	}
}

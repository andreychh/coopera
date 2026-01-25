// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package transport

import (
	"strings"
)

type MappedError struct {
	origin   error
	old, new string
}

func NewMappedError(origin error, old, new string) *MappedError {
	return &MappedError{
		origin: origin,
		old:    old,
		new:    new,
	}
}

func (e MappedError) Error() string {
	return strings.ReplaceAll(e.origin.Error(), e.old, e.new)
}

func (e MappedError) Unwrap() error {
	return e.origin
}

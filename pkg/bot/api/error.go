// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"fmt"

	"github.com/andreychh/coopera/pkg/ptr"
)

type Error struct {
	envelope Envelope
}

func NewError(envelope Envelope) *Error {
	return &Error{
		envelope: envelope,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"telegram bot api error (code %d): %s",
		ptr.ValueOrDefault(e.envelope.ErrorCode, -1),
		ptr.ValueOrDefault(e.envelope.Description, "unknown"),
	)
}

func (e *Error) Envelope() Envelope {
	return e.envelope
}

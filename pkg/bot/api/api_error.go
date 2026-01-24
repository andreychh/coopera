// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"fmt"

	"github.com/andreychh/coopera/pkg/utils"
)

type apiError struct {
	envelope Envelope
}

func (a apiError) Error() string {
	return fmt.Sprintf(
		"telegram bot api error (code %d): %s",
		utils.ValueOrDefault(a.envelope.ErrorCode, -1),
		utils.ValueOrDefault(a.envelope.Description, "unknown"),
	)
}

func (a apiError) Envelope() Envelope {
	return a.envelope
}

func NewError(envelope Envelope) Error {
	return &apiError{
		envelope: envelope,
	}
}

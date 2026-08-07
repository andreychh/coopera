// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import (
	"fmt"
	"time"
)

// DateTime is a moment as the domain speaks of it. It reads and writes
// RFC 3339, and always renders in UTC, so the same moment is spelled the
// same way wherever it is shown.
type DateTime time.Time

func ParseDateTime(s string) (DateTime, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return DateTime{}, fmt.Errorf("invalid format: %w", err)
	}
	return DateTime(t), nil
}

func (d DateTime) String() string {
	return time.Time(d).UTC().Format(time.RFC3339Nano)
}

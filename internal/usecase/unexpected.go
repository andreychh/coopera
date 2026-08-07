// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package usecase

import (
	"fmt"

	"github.com/andreychh/coopera/internal/domain"
)

// unexpected labels a failure the domain has no answer for. It only
// spells the wrapping out once; every caller still says which step broke.
func unexpected(step string, err error) domain.UnexpectedError {
	return domain.UnexpectedError{Err: fmt.Errorf("%s: %w", step, err)}
}

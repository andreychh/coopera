// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package usecase

import (
	"context"

	"github.com/andreychh/coopera/internal/domain"
)

// Post carries a sign-in code to an address. It is the one thing in the
// way of signing in that this system does not own: whether the letter
// arrives is decided by mail servers, spam filters and the person's
// provider, and none of them report back.
//
// Send returning nil means the letter was handed over, not that it was
// read or even delivered.
type Post interface {
	Send(ctx context.Context, to domain.Email, code domain.SignInCode) error
}

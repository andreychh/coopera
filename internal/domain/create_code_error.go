// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import (
	"fmt"
	"time"
)

// CreateCodeError is everything that can go wrong once asking for a
// sign-in code has begun: [TooManyCodesError] answers 429 and
// [UnexpectedError] answers 500. A malformed address is absent on
// purpose — the handler answers 400 before the usecase runs.
//
// There is no failure for "no such person" here, and there cannot be:
// the answer is the same whether or not an account stands behind the
// address, or the form would become a directory of who is registered.
//
//sumtype:decl
type CreateCodeError interface {
	error

	createCodeError()
}

// TooManyCodesError says a code was asked for too soon. Two limits raise
// it — one a minute apart, five within any sixty minutes — and they are
// not told apart, because the answer to either is the same: wait. How
// long is the only thing that differs, and RetryAfter carries it.
//
// The limits exist to keep the system from becoming an instrument: the
// letters are sent by it, and received by whoever was written into the
// form.
type TooManyCodesError struct {
	RetryAfter time.Duration
}

func (e TooManyCodesError) Error() string {
	return fmt.Sprintf("too many codes requested, retry after %s", e.RetryAfter)
}

func (TooManyCodesError) createCodeError() {}

func (UnexpectedError) createCodeError() {}

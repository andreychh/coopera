// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// CreateSessionError is everything that can go wrong once trading a code
// for a pass has begun: [SignInCodeMismatchError] answers 401,
// [SignInCodeNotUsableError] answers 410, and [UnexpectedError] answers
// 500. A malformed address or code is absent on purpose — the handler
// answers 400 before the usecase runs.
//
// There is no failure for "no such person", and there cannot be one: an
// address with nothing behind it gets an account here rather than a
// refusal. Signing in and signing up are one act, and a set of failures
// that told them apart would turn the sign-in form into the directory
// this system refuses to be.
//
// Nor is there one for "not introduced yet". Coming in without a name is
// the ordinary way to arrive — every account starts there — and the pass
// issued to a nameless person is the same pass. What it opens is decided
// later, and elsewhere.
//
//sumtype:decl
type CreateSessionError interface {
	error

	createSessionError()
}

func (SignInCodeMismatchError) createSessionError() {}

func (SignInCodeNotUsableError) createSessionError() {}

func (UnexpectedError) createSessionError() {}

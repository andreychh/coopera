// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// CreateIntroductionError is everything that can go wrong once
// introducing oneself has begun: [UserNotFoundError] answers 401, both
// [AlreadyIntroducedError] and [UsernameTakenError] answer 409, and
// [UnexpectedError] answers 500. A malformed id or name is absent on
// purpose — the handler answers 400 before the usecase runs.
//
// The two conflicts share a status because they are the same kind of
// refusal: the request is well formed, and something that already exists
// stands in its way. They differ in what exists — a name this person
// gave earlier, or a name somebody else holds — and until problems carry
// a type, that difference lives only in the text of the answer.
//
//sumtype:decl
type CreateIntroductionError interface {
	error

	createIntroductionError()
}

func (UserNotFoundError) createIntroductionError() {}

func (AlreadyIntroducedError) createIntroductionError() {}

func (UsernameTakenError) createIntroductionError() {}

func (UnexpectedError) createIntroductionError() {}

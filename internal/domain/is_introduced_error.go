// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// IsIntroducedError is everything that can go wrong once asking whether
// someone has introduced themselves: [UserNotFoundError] answers 401,
// [UnexpectedError] answers 500.
//
// A person without a name is not among these. That is the answer the
// question was asked for, not a failure of asking it — the refusal it
// leads to belongs to whatever was being guarded, and is spelled
// [NotIntroducedError] there.
//
//sumtype:decl
type IsIntroducedError interface {
	error

	isIntroducedError()
}

func (UserNotFoundError) isIntroducedError() {}

func (UnexpectedError) isIntroducedError() {}

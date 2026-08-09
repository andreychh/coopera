// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// GetMeError is everything that can go wrong once fetching the caller
// has begun: [UserNotFoundError] answers 401, [UnexpectedError] answers
// 500. A malformed caller id is absent on purpose — the handler answers
// 400 before the usecase runs.
//
// Nobody can be refused their own record: there is no owner to ask and
// nothing to be hidden from oneself. So the only way this fails on the
// caller's account is that the caller is nobody at all — which is a
// question about who is asking, not about what they asked for.
//
// That is why [UserNotFoundError] is not expected to stay. While
// X-User-Id stands in for authentication, the id in the header names
// someone who may not exist, and each usecase finds this out for itself
// at the moment it reaches the database. Once authentication settles who
// the caller is before any usecase runs, nobody unknown gets this far,
// and this set shrinks to one case.
//
//sumtype:decl
type GetMeError interface {
	error

	getMeError()
}

func (UserNotFoundError) getMeError() {}

func (UnexpectedError) getMeError() {}

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// ListMyTeamsError is everything that can go wrong once listing one's
// own teams has begun, which is only [UnexpectedError] and answers 500.
// A malformed caller id is absent on purpose — the handler answers 400
// before the usecase runs.
//
// The single case is the point, not an oversight. Asking for one's own
// teams cannot be refused: there is no owner to ask, nothing to be
// forbidden from, and having no teams is an ordinary answer rather than
// a failure. Nothing the caller can do makes this fail, so anything that
// does fail is ours.
//
// It stays a sum so that a second case, if one is ever found, breaks the
// build at every place failures become responses.
//
//sumtype:decl
type ListMyTeamsError interface {
	error

	listMyTeamsError()
}

func (UnexpectedError) listMyTeamsError() {}

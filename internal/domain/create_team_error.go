// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// CreateTeamError is everything that can go wrong once creating a team
// has begun: [UserNotFoundError] is 401, [UnexpectedError] is 500. A
// malformed name is absent on purpose — the handler answers 400 before
// the usecase runs, and by then [TeamName] guarantees the name is one
// the database will accept.
//
// So the CHECK constraint on team names is not a case here. If it ever
// fires, a name got past [ParseTeamName] that should not have, which is
// a fault in this program rather than in the request: it belongs in
// [UnexpectedError] and answers 500, never 400.
//
// Two teams may share a name, so there is no case for that either.
//
// [UserNotFoundError] is here for as long as X-User-Id stands in for
// authentication: it is the usecase that finds out the actor is unknown.
// Real authentication moves that discovery ahead of the usecase, and the
// case then leaves this set.
//
//sumtype:decl
type CreateTeamError interface {
	error

	createTeamError()
}

func (UserNotFoundError) createTeamError() {}

func (UnexpectedError) createTeamError() {}

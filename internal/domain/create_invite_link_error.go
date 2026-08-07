// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// CreateInviteLinkError is everything that can go wrong once creating an
// invite link has begun: [NotTeamOwnerError] is 403, [UnexpectedError]
// is 500. Malformed input is absent on purpose — the handler answers 400
// before the usecase runs.
//
// There is no case for a missing team, deliberately. An outsider must
// not learn whether a team exists, so "no such team", "not a member" and
// "a member but not the owner" are one indistinguishable refusal. The
// absent case is what enforces it: the handler cannot answer 404,
// because it has nothing to answer 404 for.
//
// The marker method is unexported, so the set is closed, and
// gochecksumtype checks that type switches cover it: a third failure
// breaks the build everywhere failures become responses, instead of
// silently falling through to 500.
//
//sumtype:decl
type CreateInviteLinkError interface {
	error

	createInviteLinkError()
}

func (NotTeamOwnerError) createInviteLinkError() {}

// UnexpectedError is a failure the domain has no answer for: the
// database being unreachable, a query breaking, a row the schema should
// have ruled out. It carries the cause for logging, and callers turn it
// into a 500 without inspecting it.
type UnexpectedError struct {
	Err error
}

func (e UnexpectedError) Error() string {
	return e.Err.Error()
}

// Unwrap exposes the cause to [errors.Is] and [errors.As].
func (e UnexpectedError) Unwrap() error {
	return e.Err
}

func (UnexpectedError) createInviteLinkError() {}

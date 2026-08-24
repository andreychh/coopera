// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// AcceptInviteLinkError is everything that can go wrong once joining a
// team by invitation has begun: [InviteLinkNotFoundError] is 404,
// [AlreadyMemberError] is 409, [InviteLinkNotUsableError] is 410,
// [UserNotFoundError] is 401, [UnexpectedError] is 500. Malformed input
// is absent on purpose — the handler answers 400 before the usecase
// runs.
//
// [AlreadyMemberError] is a failure rather than an outcome. Joining
// happens once, and reporting a second attempt as success would make one
// success code stand for two different things — a membership that came
// into being and one that was already there — with nothing in the answer
// to tell the caller which they got. Naming it a conflict is the same
// reading introducing yourself twice already gets.
//
// Revocation and expiry are one case, not two. Either way the person
// needs a fresh invitation, so the difference cannot change what they
// do; telling them which one it was only reports a decision somebody
// else made. A missing link stays its own case, though — codes cannot be
// guessed, so saying one doesn't exist gives nothing away, and it lets a
// mistyped code be recognised as one.
//
// [UserNotFoundError] is here for as long as X-User-Id stands in for
// authentication: it is the usecase that finds out the actor is unknown.
// Real authentication moves that discovery ahead of the usecase, and the
// case then leaves this set.
//
//sumtype:decl
type AcceptInviteLinkError interface {
	error

	acceptInviteLinkError()
}

func (InviteLinkNotFoundError) acceptInviteLinkError() {}

func (AlreadyMemberError) acceptInviteLinkError() {}

func (InviteLinkNotUsableError) acceptInviteLinkError() {}

func (UserNotFoundError) acceptInviteLinkError() {}

func (UnexpectedError) acceptInviteLinkError() {}

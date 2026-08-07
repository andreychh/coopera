// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// RevokeInviteLinkError is everything that can go wrong once revoking an
// invite link has begun: [InviteLinkNotFoundError] is 404,
// [NotTeamOwnerError] is 403, [InviteLinkAlreadyRevokedError] is 409,
// [UnexpectedError] is 500. Malformed input is absent on purpose — the
// handler answers 400 before the usecase runs.
//
// Unlike creating or listing links, this set does distinguish a missing
// team from an unowned one — or rather, it never has to. Revoking is
// reached by code, and a code cannot be guessed, so whoever gets this
// far already knows the link exists and with it the team. There is
// nothing left to conceal.
//
// Only an already revoked link is a conflict. A link that merely ran out
// of time can still be revoked: expiring is something that happened,
// revoking is something the owner means, and meaning it about a link
// that has lapsed is not a mistake.
//
//sumtype:decl
type RevokeInviteLinkError interface {
	error

	revokeInviteLinkError()
}

func (InviteLinkNotFoundError) revokeInviteLinkError() {}

func (NotTeamOwnerError) revokeInviteLinkError() {}

func (InviteLinkAlreadyRevokedError) revokeInviteLinkError() {}

func (UnexpectedError) revokeInviteLinkError() {}

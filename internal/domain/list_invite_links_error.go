// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// ListInviteLinksError is everything that can go wrong once listing a
// team's invite links has begun: [NotTeamOwnerError] is 403,
// [UnexpectedError] is 500. Malformed input is absent on purpose — the
// handler answers 400 before the usecase runs.
//
// As with creating a link, there is no case for a missing team: an
// outsider must not learn that a team exists, so "no such team", "not a
// member" and "a member but not the owner" are one indistinguishable
// refusal, and the absent case is what enforces it.
//
// An owner whose team has no links is not a failure at all. That is an
// empty list, and keeping it apart from a refusal is why listing asks
// about ownership separately instead of reading emptiness as an answer.
//
//sumtype:decl
type ListInviteLinksError interface {
	error

	listInviteLinksError()
}

func (NotTeamOwnerError) listInviteLinksError() {}

func (UnexpectedError) listInviteLinksError() {}

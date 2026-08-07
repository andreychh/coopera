// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// GetTeamError is everything that can go wrong once fetching a team has
// begun: [TeamNotFoundError] is 404, [UnexpectedError] is 500. A
// malformed id is absent on purpose — the handler answers 400 before the
// usecase runs.
//
// [TeamNotFoundError] covers three things at once: no such team, the
// actor never belonged to it, and the actor has left it. An outsider
// must not learn that a team exists, and here that refusal reads as
// absence rather than as prohibition — there is nothing to forbid
// someone from, if for them the team is not there.
//
// Acting on a team refuses differently, with 403, and that is not a
// contradiction. Those refusals can reach a member, someone who does see
// the team; answering them with 404 would be a lie they could catch by
// fetching the team itself. Only a non-member is ever refused here, and
// for them absence is the whole truth.
//
//sumtype:decl
type GetTeamError interface {
	error

	getTeamError()
}

func (TeamNotFoundError) getTeamError() {}

func (UnexpectedError) getTeamError() {}

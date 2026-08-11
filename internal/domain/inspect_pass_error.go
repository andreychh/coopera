// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// InspectPassError is everything that can go wrong once the door has
// begun looking at a pass: [PassNotUsableError] answers 401 and
// [UnexpectedError] answers 500. A token that cannot be read at all is
// absent on purpose — the door refuses it before asking anything.
//
// A person without a name is not among these either. That is the answer
// the question was asked for, not a failure of asking it, and the
// refusal it leads to belongs to whatever was being guarded.
//
//sumtype:decl
type InspectPassError interface {
	error

	inspectPassError()
}

func (PassNotUsableError) inspectPassError() {}

func (UnexpectedError) inspectPassError() {}

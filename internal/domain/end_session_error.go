// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// EndSessionError is everything that can go wrong once leaving has
// begun, which is only [UnexpectedError], answering 500.
//
// The set is this short because leaving is the one act nobody may be
// refused. A pass that cannot be read is turned away at the door before
// this is reached, and every other outcome is the one asked for: a
// session already closed leaves nothing to close, and asking again is
// not a mistake but a wish already granted.
//
// It stays a sum of one rather than a plain error so that a refusal
// added later has to be answered for at the handler, where
// gochecksumtype will notice it.
//
//sumtype:decl
type EndSessionError interface {
	error

	endSessionError()
}

func (UnexpectedError) endSessionError() {}

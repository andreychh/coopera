// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// RefreshSessionError is everything that can go wrong once renewing a
// pass has begun: [RefreshTokenNotUsableError] answers 401 and
// [UnexpectedError] answers 500.
//
// One refusal covers every way a key can fail, and that is the whole
// point rather than a shortcut — see [RefreshTokenNotUsableError]. A
// theft caught in the act is not a separate case here either: it ends
// the session as it goes, and what the one holding the key is told
// afterwards is the same as for the other three.
//
// There is nothing here about the person being unknown or unnamed. The
// key names its session and the session names its person; a nameless
// person renews their pass like anyone else, since renewing shows them
// to nobody.
//
//sumtype:decl
type RefreshSessionError interface {
	error

	refreshSessionError()
}

func (RefreshTokenNotUsableError) refreshSessionError() {}

func (UnexpectedError) refreshSessionError() {}

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// Actor is whoever the request came from, as the pass they showed
// describes them: which person, and which of their sessions they are
// speaking from.
//
// Both are needed and neither replaces the other. The person is who
// everything is done on behalf of; the session is what "here" means in
// "sign out here", and what a list of devices would name. One person has
// several sessions at once — a phone and a laptop are the ordinary case,
// not an oddity to be prevented.
type Actor struct {
	ID        ID
	SessionID ID
}

// PassNotUsableError says the pass shown is no longer good.
//
// Three ways lead here: nobody by that name exists, the session named
// alongside does not, or it has been closed — by signing out or by a
// theft being noticed. They are one refusal on purpose. All three leave
// the holder with the same single thing to do, and telling them apart
// would report to whoever holds a stolen pass exactly what became of it.
type PassNotUsableError struct{}

func (PassNotUsableError) Error() string {
	return "pass is not usable"
}

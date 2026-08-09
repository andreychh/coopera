// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// UserState is how far a person has got towards being someone others
// can see. It has exactly two cases, [Unintroduced] and [Introduced]:
// the marker method is unexported, so no package can add a third, and
// gochecksumtype checks that type switches cover both.
//
// A state is derived, never asserted — build it with [NewUserState] so
// that every caller gets the same answer from the same facts.
//
//sumtype:decl
type UserState interface {
	userState()
}

// NewUserState works out where a person stands from the one fact that
// settles it: whether they have a username. A nil username means they
// have never introduced themselves.
//
// The database records the moment of introduction separately, and a
// CHECK keeps the two in step, so either would answer this question. The
// username answers it and is needed anyway.
func NewUserState(username *Username) UserState {
	if username == nil {
		return Unintroduced{}
	}
	return Introduced{Username: *username}
}

// Unintroduced is the state of someone the system knows of but cannot
// show to anyone. It carries no name — not an empty one, none at all,
// so "introduced without a name" cannot be written down.
type Unintroduced struct{}

func (Unintroduced) userState() {}

// Introduced is the state of someone who has named themselves. It
// carries that name, which is also how everyone else tells them apart.
type Introduced struct {
	Username Username
}

func (Introduced) userState() {}

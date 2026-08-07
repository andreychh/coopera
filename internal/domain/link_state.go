// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import "time"

// LinkState is the condition an invite link is in. It has exactly three
// cases, [Active], [Expired] and [Revoked]: the marker method is
// unexported, so no package can add a fourth, and gochecksumtype checks
// that type switches cover all three.
//
// A state is derived, never asserted — build it with [NewLinkState] so
// that every caller gets the same answer from the same facts.
//
//sumtype:decl
type LinkState interface {
	linkState()
}

// NewLinkState works out the state of a link from what is stored about
// it: when it was set to expire, when it was revoked, and what time it
// is now. A nil expiresAt means the link was issued without a deadline,
// a nil revokedAt means it was never revoked.
//
// Revocation wins over expiry: a link that was revoked and has since
// passed its deadline reads as [Revoked], because revoking is something
// the owner did and expiring is something that merely happened.
//
// A link expires the moment its deadline arrives, not a moment after.
func NewLinkState(expiresAt, revokedAt *DateTime, now DateTime) LinkState {
	if revokedAt != nil {
		return Revoked{At: *revokedAt}
	}
	if expiresAt == nil {
		return Active{Validity: Forever{}}
	}
	if !time.Time(*expiresAt).After(time.Time(now)) {
		return Expired{At: *expiresAt}
	}
	return Active{Validity: Until{Time: *expiresAt}}
}

// Active is the state of a link that still lets people in, for as long
// as its [Validity] says.
type Active struct {
	Validity Validity
}

func (Active) linkState() {}

// Expired is the state of a link whose deadline has passed. It carries
// that deadline, and having one is what makes the case possible at all:
// a link issued as [Forever] can never reach it.
type Expired struct {
	At DateTime
}

func (Expired) linkState() {}

// Revoked is the state of a link the owner took back. It carries the
// moment that happened.
type Revoked struct {
	At DateTime
}

func (Revoked) linkState() {}

// LinkStatus names one case of [LinkState]. It exists to pick links out
// by their state, not to represent that state: a LinkStatus is what a
// caller is asking for, and the only thing it can answer is whether a
// given link is it.
type LinkStatus string

const (
	StatusActive  LinkStatus = "active"
	StatusExpired LinkStatus = "expired"
	StatusRevoked LinkStatus = "revoked"
)

// Matches reports whether state is the case this status names. An
// unrecognised status matches nothing, so a filter that means nothing
// selects nothing rather than everything.
func (s LinkStatus) Matches(state LinkState) bool {
	switch state.(type) {
	case Active:
		return s == StatusActive
	case Expired:
		return s == StatusExpired
	case Revoked:
		return s == StatusRevoked
	}
	return false
}

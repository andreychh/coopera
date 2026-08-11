// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package usecase

import (
	"time"

	"github.com/andreychh/coopera/internal/domain"
)

// Seal turns a person and a session into an access token, and is the one
// thing in the system holding the key that token is signed with. Nothing
// else can make one, and nothing else needs to.
//
// The lifetime comes back alongside the token because the deadline is
// written inside it, by this. Naming the same span in a second place is
// how the two would eventually disagree.
//
// Only stamping is declared here. Reading a token back needs the same
// key but happens elsewhere, at the door, and what is needed there will
// be said there.
type Seal interface {
	Stamp(
		personID domain.ID,
		sessionID domain.ID,
	) (domain.AccessToken, time.Duration, error)
}

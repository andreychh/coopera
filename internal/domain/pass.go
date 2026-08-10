// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
)

// AccessToken is what a person shows on every request. It is signed
// rather than stored: what it says about its holder is checked against a
// key, so nothing has to be looked up to learn who is asking.
//
// The price of not storing it is not being able to take it back. It
// therefore lives minutes, and anything that must stop sooner than that
// stops through the session behind it.
type AccessToken string

func (t AccessToken) String() string {
	return string(t)
}

// RefreshToken buys a fresh pair once the access token dies. It is
// opaque — random bytes with nothing readable inside — because it is
// stored, and a secret that is looked up has no need to carry its own
// meaning.
//
// It is the one worth stealing: it outlives the access token many times
// over and buys new ones. Spending it burns it, and a second use of the
// same one is how a loose copy gives itself away.
type RefreshToken string

// NewRefreshToken draws a token from the system's source of randomness.
// Thirty-two bytes are past guessing by any margin that matters, and
// nothing about its holder is woven in: the token stands for a row and
// says nothing on its own, so learning one teaches nothing about the
// next.
func NewRefreshToken() (RefreshToken, error) {
	var b [32]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return "", fmt.Errorf("draw random bytes: %w", err)
	}
	return RefreshToken(base64.RawURLEncoding.EncodeToString(b[:])), nil
}

func (t RefreshToken) String() string {
	return string(t)
}

// Hash is what the system keeps in place of the token, so that a copy of
// the table opens nothing.
//
// A plain SHA-256, with no salt and no work factor, and both absences
// are deliberate. Those exist to slow down guessing at something a
// person chose, out of a set small enough to walk through. This is
// thirty-two random bytes: there is no set to walk, and nothing to slow
// down.
func (t RefreshToken) Hash() string {
	sum := sha256.Sum256([]byte(t))
	return hex.EncodeToString(sum[:])
}

// Pass is the pair of keys to a session, and not the session itself: a
// session lasts from signing in until signing out, while these are
// replaced many times over its life.
//
// Both lifetimes are spans rather than moments because whoever holds
// them counts from receiving them. Naming instants would send the holder
// to its own clock, and a phone kept off the network drifts by longer
// than an access token lives.
type Pass struct {
	AccessToken      AccessToken
	AccessExpiresIn  time.Duration
	RefreshToken     RefreshToken
	RefreshExpiresIn time.Duration
}

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// SignInCode is the six digits mailed to an address, and the whole of what a
// person shows to be let in. It means nothing on its own: the same six
// digits stand for different codes at different addresses, so a code is
// only ever checked together with the address it was issued for.
//
// Six digits are short enough to carry from a mail client to another
// device by eye, and that is the reason for the length. What makes them
// hard to guess is not the number of them but what stands around it: ten
// minutes of life, five attempts, and one address.
type SignInCode string

// NewSignInCode draws six digits from the system's source of randomness.
//
// It is crypto/rand rather than math/rand because a code is a secret,
// and a sequence anyone can continue after seeing a few of them is not
// one. Leading zeros are kept — "012345" is as good a code as any, and
// dropping them would quietly shrink the range.
func NewSignInCode() (SignInCode, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", fmt.Errorf("draw random number: %w", err)
	}
	return SignInCode(fmt.Sprintf("%06d", n)), nil
}

func (c SignInCode) String() string {
	return string(c)
}

// CodeDelivery is what the person asking for a code is told: how long
// the one just sent keeps working, and how long before another may be
// asked for. The code itself is not here — it belongs in the mailbox.
//
// Both are spans rather than moments because nothing here is stored.
// This is a thing said to someone, and the one it is said to counts from
// hearing it: a countdown on a screen, a timer before the button comes
// back. Naming instants would send them to a clock the system does not
// own and has no reason to trust.
type CodeDelivery struct {
	ExpiresIn  time.Duration
	RetryAfter time.Duration
}

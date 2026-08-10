// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
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

// ParseSignInCode accepts six digits and refuses everything else.
//
// Refusing here rather than letting the characters reach the table is
// not tidiness. A code that could never have been issued is not a wrong
// guess, and charging an attempt for one would hand anybody a way to
// burn a stranger's code with nonsense: five requests of "hello" and the
// letter in their mailbox stops working.
func ParseSignInCode(s string) (SignInCode, error) {
	if !signInCodeShape.MatchString(s) {
		return "", errors.New("must be six digits")
	}
	return SignInCode(s), nil
}

// signInCodeShape is the rule the schema puts on the column, written in
// the notation Go prefers: what the table refuses to hold is refused
// before it ever gets there.
var signInCodeShape = regexp.MustCompile(`^\d{6}$`)

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

// SignInCodeMismatchError says the digits do not match the ones issued
// for this address. The code survives it: attempts are there for a slip
// of the finger, and only the fifth wrong one burns it.
//
// AttemptsLeft is told to the person because it is the only thing here
// that changes what they do next — look at the letter again, or ask for
// a new code.
type SignInCodeMismatchError struct {
	AttemptsLeft int64
}

func (e SignInCodeMismatchError) Error() string {
	return fmt.Sprintf("sign-in code does not match, %d attempts left", e.AttemptsLeft)
}

// SignInCodeNotUsableError says no live code stands behind the address.
// Five ways lead here — it expired, it was already spent, five wrong
// attempts burned it, a newer code replaced it, or none was ever asked
// for — and they are one refusal on purpose.
//
// Partly because the remedy is the same for all five: ask for a code.
// But mostly because whoever typed the address is not always its owner,
// and two of the five would report the owner's doings to them. "Already
// spent" says somebody got in moments ago; "replaced by a newer one"
// says somebody asked again. Neither is the typist's business.
type SignInCodeNotUsableError struct {
	Email Email
}

func (e SignInCodeNotUsableError) Error() string {
	return fmt.Sprintf("no usable sign-in code for %s", e.Email)
}

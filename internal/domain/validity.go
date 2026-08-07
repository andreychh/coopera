// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import (
	"errors"
	"time"
)

// Validity is how long an invite link stays usable. It has exactly two
// cases, [Until] and [Forever]: the marker method is unexported, so no
// package can add a third, and gochecksumtype checks that type switches
// cover both.
//
// Its zero value is a nil interface, which is neither case. A Validity
// must come from [ParseValidity] or be [Forever], and switches still
// need a default clause to reject the zero value.
//
//sumtype:decl
type Validity interface {
	validity()
}

// ParseValidity reads a deadline in RFC 3339 format, failing if it isn't
// a well-formed timestamp or doesn't lie in the future. A link that
// never expires is [Forever] rather than a parse result: there is no
// text to read.
func ParseValidity(s string) (Validity, error) {
	at, err := ParseDateTime(s)
	if err != nil {
		return nil, err
	}
	if !time.Time(at).After(time.Now()) {
		return nil, errors.New("must be in the future")
	}
	return Until{Time: at}, nil
}

// Until is the validity of a link that stops working at a fixed moment.
type Until struct {
	Time DateTime
}

func (Until) validity() {}

// Forever is the validity of a link that never stops working.
type Forever struct{}

func (Forever) validity() {}

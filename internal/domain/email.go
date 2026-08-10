// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

// emailShape is the pattern the HTML standard defines for
// <input type="email">, which calls itself a wilful violation of RFC
// 5322 because the real grammar allows addresses no mail system on the
// web has ever seen. See https://html.spec.whatwg.org/#valid-e-mail-address.
//
// Two departures, both towards being stricter: the domain must carry a
// dot, or the address is one nothing can be delivered to; and the local
// part is a dot-atom, so a dot may not lead, trail or double, which the
// standard permits and RFC 5322 does not. The same pattern stands in the
// schema, where it is the guarantee; here it is what lets a person be
// told what is wrong.
var emailShape = regexp.MustCompile(
	"^[a-z0-9!#$%&'*+/=?^_`{|}~-]+" +
		"(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@" +
		"[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?" +
		"(?:\\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)+$",
)

// Email is the address a person is known by and the only way to reach
// them. Unlike [Username] it is never shown to anyone else.
type Email string

// ParseEmail accepts what can plausibly receive mail: at most 254
// characters overall, at most 64 before the "@", lower case throughout,
// and the shape described by [emailShape].
//
// Capitals are refused rather than folded, as in [ParseUsername].
// Folding would accept what the schema calls malformed and answer 202
// where the contract promises 400; refusing keeps one address to one
// spelling, and leaves the client to lowercase what someone typed before
// sending it. Two addresses differing only in case are the same address,
// and this is how that is enforced — by admitting one of the two.
//
// Nothing else is folded: "+tags" and dots inside the local part stay,
// because whether they distinguish two mailboxes is for the receiving
// server to say, not for this one.
//
// Whether mail actually arrives is not knowable here. Only delivery
// tells that, and its outcome is deliberately never reported.
func ParseEmail(s string) (Email, error) {
	if utf8.RuneCountInString(s) > 254 {
		return "", errors.New("must be at most 254 characters")
	}
	if s != strings.ToLower(s) {
		return "", errors.New("must be lowercase")
	}
	local, _, found := strings.Cut(s, "@")
	if !found {
		return "", errors.New("must contain an @")
	}
	if utf8.RuneCountInString(local) > 64 {
		return "", errors.New("must have at most 64 characters before the @")
	}
	if !emailShape.MatchString(s) {
		return "", errors.New("must be a well-formed email address")
	}
	return Email(s), nil
}

func (e Email) String() string {
	return string(e)
}

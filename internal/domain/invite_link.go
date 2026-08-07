// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// Code is a high-entropy invite link credential: knowing it is what grants
// access, not just a lookup key, so it's generated with crypto/rand rather
// than a predictable id.
type Code string

func NewCode() (Code, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	return Code(base64.RawURLEncoding.EncodeToString(b)), nil
}

func (c Code) String() string {
	return string(c)
}

// InviteLink is a link as it is reported to whoever asked for it. Its
// deadline and its revocation live inside [LinkState] rather than beside
// it, so a revoked link always carries the moment it was revoked and an
// expired one always carries its deadline.
type InviteLink struct {
	Code      Code
	UseCount  int64
	State     LinkState
	CreatedAt DateTime
}

// InviteLinkNotFoundError says no link answers to this code. It is kept
// apart from InviteLinkNotUsableError deliberately: a code cannot be
// guessed, so admitting that one does not exist gives nothing away, and
// it lets a mistyped code be recognised as one.
type InviteLinkNotFoundError struct {
	Code Code
}

func (e InviteLinkNotFoundError) Error() string {
	return fmt.Sprintf("invite link %s not found", e.Code)
}

// InviteLinkNotUsableError says the link no longer lets anyone in,
// whether it was revoked or ran out of time. The two are one refusal on
// purpose: either way a fresh invitation is needed, and which of them it
// was is the owner's to explain.
type InviteLinkNotUsableError struct {
	Code Code
}

func (e InviteLinkNotUsableError) Error() string {
	return fmt.Sprintf("invite link %s is expired or revoked", e.Code)
}

// InviteLinkAlreadyRevokedError says there is nothing left to revoke. A
// link that merely lapsed is not this: expiring is something that
// happened, revoking is something the owner means.
type InviteLinkAlreadyRevokedError struct {
	Code Code
}

func (e InviteLinkAlreadyRevokedError) Error() string {
	return fmt.Sprintf("invite link %s is already revoked", e.Code)
}

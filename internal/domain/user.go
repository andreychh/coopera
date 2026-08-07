// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

// ForeignKeyViolation is the Postgres SQLSTATE code for a violated
// foreign key (https://www.postgresql.org/docs/current/errcodes-appendix.html),
// and MembersUserIDForeignKey is the name Postgres generates for
// members' unnamed "user_id UUID NOT NULL REFERENCES users (id)"
// constraint, following its default <table>_<column>_fkey convention.
//
// This is the last place where the domain speaks Postgres, and it is a
// brittle one: renaming the constraint in a migration would break the
// match here without a word from the compiler. It survives only because
// nothing but the database can tell that an actor does not exist. Every
// other such translation has been replaced by a query whose result
// carries the answer.
const (
	ForeignKeyViolation     = "23503"
	MembersUserIDForeignKey = "members_user_id_fkey"
)

type ID uuid.UUID

func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}

func (i ID) String() string {
	return uuid.UUID(i).String()
}

type TeamName string

func ParseTeamName(s string) (TeamName, error) {
	if strings.TrimSpace(s) != s {
		return "", errors.New("must not have leading or trailing whitespace")
	}
	count := utf8.RuneCountInString(s)
	if count < 1 || count > 100 {
		return "", errors.New("must be between 1 and 100 characters")
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return "", errors.New("must not contain control characters")
		}
	}
	return TeamName(s), nil
}

func (n TeamName) String() string {
	return string(n)
}

type DateTime time.Time

func ParseDateTime(s string) (DateTime, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return DateTime{}, fmt.Errorf("invalid format: %w", err)
	}
	return DateTime(t), nil
}

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

func (d DateTime) String() string {
	return time.Time(d).UTC().Format(time.RFC3339Nano)
}

type UserNotFoundError struct {
	ID ID
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user %s not found", e.ID)
}

type TeamNotFoundError struct {
	ID ID
}

func (e TeamNotFoundError) Error() string {
	return fmt.Sprintf("team %s not found", e.ID)
}

type NotTeamOwnerError struct {
	TeamID ID
}

func (e NotTeamOwnerError) Error() string {
	return fmt.Sprintf("caller is not owner of team %s", e.TeamID)
}

type InviteLinkNotFoundError struct {
	Code Code
}

func (e InviteLinkNotFoundError) Error() string {
	return fmt.Sprintf("invite link %s not found", e.Code)
}

type InviteLinkNotUsableError struct {
	Code Code
}

func (e InviteLinkNotUsableError) Error() string {
	return fmt.Sprintf("invite link %s is expired or revoked", e.Code)
}

type InviteLinkAlreadyRevokedError struct {
	Code Code
}

func (e InviteLinkAlreadyRevokedError) Error() string {
	return fmt.Sprintf("invite link %s is already revoked", e.Code)
}

type Team struct {
	ID        ID
	Name      TeamName
	CreatedAt DateTime
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

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Username is what one person is called and how everyone else tells
// them apart. Unlike [TeamName] it identifies: no two people may share
// one.
type Username string

// ParseUsername accepts three to thirty-two characters — lowercase latin
// letters, digits and underscores — with no underscore at either end and
// never two in a row.
//
// The alphabet is narrow because the name identifies. A name that
// resembles another can then be used to impersonate: "anya" in Latin and
// "аnya" with a Cyrillic "а" are one image and two people. Within a
// single script there is no such pair.
//
// Capitals are refused rather than folded. Folding would accept what the
// schema calls malformed and answer 201 where the contract promises 400;
// refusing keeps one person to one spelling, and leaves the client to
// lowercase what someone typed before sending it.
func ParseUsername(s string) (Username, error) {
	count := utf8.RuneCountInString(s)
	if count < 3 || count > 32 {
		return "", errors.New("must be between 3 and 32 characters")
	}
	for _, r := range s {
		letter := r >= 'a' && r <= 'z'
		digit := r >= '0' && r <= '9'
		if !letter && !digit && r != '_' {
			return "", errors.New(
				"must contain only lowercase latin letters, digits and underscores",
			)
		}
	}
	if strings.HasPrefix(s, "_") || strings.HasSuffix(s, "_") {
		return "", errors.New("must not start or end with an underscore")
	}
	if strings.Contains(s, "__") {
		return "", errors.New("must not contain two underscores in a row")
	}
	return Username(s), nil
}

func (n Username) String() string {
	return string(n)
}

// User is a person as they are reported to whoever asked about them.
// Whether they have a name at all is carried by [UserState], so someone
// nameless cannot be mistaken for someone named an empty string.
type User struct {
	ID        ID
	State     UserState
	CreatedAt DateTime
}

// UserNotFoundError says the actor of a request is nobody the system
// knows. While X-User-Id stands in for authentication, each usecase
// discovers this for itself; real authentication settles it before a
// usecase runs, and this error leaves the usecases' failure sets.
type UserNotFoundError struct {
	ID ID
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user %s not found", e.ID)
}

// AlreadyIntroducedError says the person has named themselves before.
// Introducing oneself happens once; renaming is a different act, with
// its own entry in the journal, and the system does not offer it yet.
type AlreadyIntroducedError struct {
	ID ID
}

func (e AlreadyIntroducedError) Error() string {
	return fmt.Sprintf("user %s is already introduced", e.ID)
}

// UsernameTakenError says the name belongs to somebody else.
//
// This does reveal that a person with that name exists, and that is
// unavoidable: a name that identifies cannot be claimed twice, so anyone
// choosing one has to be told when their choice is gone. Nothing else
// leaks with it — the answer names no person, only an occupied name.
type UsernameTakenError struct {
	Username Username
}

func (e UsernameTakenError) Error() string {
	return fmt.Sprintf("username %s is taken", e.Username)
}

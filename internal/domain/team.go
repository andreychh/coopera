// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// TeamName is what people call a team among themselves, not how the
// system tells one team from another: two teams may share a name, and
// the ID is what distinguishes them.
type TeamName string

// ParseTeamName accepts a name of one to a hundred characters — counted
// as characters and not as bytes, so a name in Cyrillic may be as long
// as one in Latin — with no control characters and no space at either
// end. A name must read as it looks: two names differing invisibly are
// one name.
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

// Team is a team as it is reported to whoever asked for it.
type Team struct {
	ID        ID
	Name      TeamName
	CreatedAt DateTime
}

// TeamNotFoundError says a team is not there for this actor, which
// covers the team not existing, the actor never having belonged to it,
// and the actor having left. The three are not told apart on purpose:
// nobody outside a team may learn that it exists.
type TeamNotFoundError struct {
	ID ID
}

func (e TeamNotFoundError) Error() string {
	return fmt.Sprintf("team %s not found", e.ID)
}

// NotTeamOwnerError says the actor may not act as owner of a team. Like
// TeamNotFoundError it covers a missing team as well as an unowned one,
// for the same reason; the two differ in what the refusal reads as, not
// in what it hides.
type NotTeamOwnerError struct {
	TeamID ID
}

func (e NotTeamOwnerError) Error() string {
	return fmt.Sprintf("caller is not owner of team %s", e.TeamID)
}

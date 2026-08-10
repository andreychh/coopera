// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package seal stamps and reads the tokens people show at the door.
package seal

import (
	"fmt"
	"time"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

// accessLifetime is how long a stamped token is accepted. It is short
// because a stamped token cannot be recalled: nothing is written when it
// is made, so nothing can be crossed out to stop it. Everything that has
// to stop sooner stops through the session behind it, and the gap
// between the two is these fifteen minutes.
const accessLifetime = 15 * time.Minute

// JWT signs what it says with a key it keeps, so that whoever reads a
// token back can trust it without looking anything up. The claims are
// plainly readable to anyone holding the token — that is what signing
// is, not hiding — so nothing goes in that its holder should not see.
//
// HS256: one key, both stamping and reading, because both happen in this
// same process. The day a second service has to read tokens without
// being able to make them, this becomes a key pair and nothing above has
// to know.
type JWT struct {
	key []byte
}

func NewJWT(key []byte) JWT {
	return JWT{key: key}
}

// claims is what a stamped token says. The session goes in sid, the name
// OpenID Connect gives it, and not in jti: jti stands for the token
// itself and has to differ between two tokens, while every token a
// session is refreshed into names that same session.
type claims struct {
	jwt.RegisteredClaims

	SessionID string `json:"sid"`
}

// Stamp writes who is asking and which session they are asking in, and
// nothing else.
//
// Whether the person has introduced themselves is deliberately absent,
// though the door asks about it on every request and could have been
// spared a lookup. It is the one fact here that changes, it changes
// exactly once, and that change is what the door exists to notice: a
// token saying "no name" would go on saying it for fifteen minutes after
// the person picked one, turning away the only caller who had just done
// what was asked of them.
func (s JWT) Stamp(
	personID domain.ID,
	sessionID domain.ID,
) (domain.AccessToken, time.Duration, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		SessionID: sessionID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   personID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessLifetime)),
		},
	})

	signed, err := token.SignedString(s.key)
	if err != nil {
		return "", 0, fmt.Errorf("sign access token: %w", err)
	}
	return domain.AccessToken(signed), accessLifetime, nil
}

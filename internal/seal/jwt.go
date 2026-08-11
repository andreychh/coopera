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

// Read says who a token names, and refuses it if the key does not
// vouch for it or its hour has passed.
//
// The order is the whole of it: the signature is checked before a single
// claim is believed. A token's claims travel in the open and anyone
// holding one can rewrite them, so read first and verify after is not a
// stricter or looser way of doing this — it is no check at all.
//
// The algorithm is pinned rather than taken from the token. A token
// carries the name of its own algorithm, and a reader that obeys it can
// be handed one saying "none", or one signed with a key it was meant to
// verify against rather than sign with. Saying in advance which
// algorithm is acceptable closes that whole family at once.
//
// The deadline is compared against this machine's clock, never the
// caller's. What the caller believes the time to be is not asked and
// would not be worth having.
func (s JWT) Read(token domain.AccessToken) (domain.Actor, error) {
	var read claims
	_, err := jwt.ParseWithClaims(
		string(token),
		&read,
		func(*jwt.Token) (any, error) { return s.key, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return domain.Actor{}, fmt.Errorf("parse access token: %w", err)
	}

	actorID, err := domain.ParseID(read.Subject)
	if err != nil {
		return domain.Actor{}, fmt.Errorf("read subject: %w", err)
	}
	sessionID, err := domain.ParseID(read.SessionID)
	if err != nil {
		return domain.Actor{}, fmt.Errorf("read session: %w", err)
	}
	return domain.Actor{ID: actorID, SessionID: sessionID}, nil
}

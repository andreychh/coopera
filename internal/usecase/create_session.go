// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/andreychh/coopera/internal/db"
	"github.com/andreychh/coopera/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateSession trades a code for a pass, and brings the account into
// being if the address had none. Signing in and signing up are one act
// here: to separate them the system would first have to say whether it
// knows the address, and that is the one thing it will not say.
//
// The code is checked together with the address it was issued for. Six
// digits are not unique — the same ones stand for different codes at
// different mailboxes — so on their own they identify nobody.
type CreateSession struct {
	pool *pgxpool.Pool
	seal Seal

	email domain.Email
	code  domain.SignInCode
}

func NewCreateSession(
	pool *pgxpool.Pool,
	seal Seal,
	email domain.Email,
	code domain.SignInCode,
) CreateSession {
	return CreateSession{pool: pool, seal: seal, email: email, code: code}
}

// Exec spends the code, finds or makes the person behind the address,
// opens a session and hands back its first pair of keys.
//
// All of it is one transaction, and it has to be: a code spent without a
// session behind it would leave its owner holding a letter that no
// longer works and an account they cannot reach.
//
// The token is stamped before the commit rather than after. Stamping can
// fail, and failing inside the transaction undoes the spending along
// with everything else — the letter still works, and the person simply
// tries again. After the commit there would be nothing left to undo.
func (u CreateSession) Exec(ctx context.Context) (domain.Pass, domain.CreateSessionError) {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return domain.Pass{}, unexpected("begin transaction", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	presented, err := db.New(tx).PresentSignInCode(ctx, db.PresentSignInCodeParams{
		Email: u.email.String(),
		Code:  u.code.String(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Pass{}, domain.SignInCodeNotUsableError{Email: u.email}
		}
		return domain.Pass{}, unexpected("present sign-in code", err)
	}
	if !presented.Matched {
		return domain.Pass{}, u.chargeAttempt(ctx, tx, presented.AttemptsLeft)
	}

	personID, err := db.New(tx).InsertUser(ctx, u.email.String())
	if err != nil {
		return domain.Pass{}, unexpected("insert user", err)
	}

	sessionID, err := db.New(tx).InsertSession(ctx, personID)
	if err != nil {
		return domain.Pass{}, unexpected("insert session", err)
	}

	refresh, err := domain.NewRefreshToken()
	if err != nil {
		return domain.Pass{}, unexpected("draw refresh token", err)
	}

	refreshExpiresIn, err := db.New(tx).InsertRefreshToken(ctx, db.InsertRefreshTokenParams{
		SessionID: sessionID,
		Hash:      refresh.Hash(),
	})
	if err != nil {
		return domain.Pass{}, unexpected("insert refresh token", err)
	}

	access, accessExpiresIn, err := u.seal.Stamp(domain.ID(personID), domain.ID(sessionID))
	if err != nil {
		return domain.Pass{}, unexpected("stamp access token", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.Pass{}, unexpected("commit", err)
	}

	return domain.Pass{
		AccessToken:      access,
		AccessExpiresIn:  accessExpiresIn,
		RefreshToken:     refresh,
		RefreshExpiresIn: time.Duration(refreshExpiresIn) * time.Second,
	}, nil
}

// chargeAttempt keeps the wrong guess on record and turns it into the
// refusal.
//
// It commits, which is the one thing about this path worth stopping at.
// Every other way out of Exec rolls back, and left to do the same, this
// one would undo the attempt it just cost — five tries would become as
// many as anybody cared to make, and the limit would exist only in the
// comments.
func (u CreateSession) chargeAttempt(
	ctx context.Context,
	tx pgx.Tx,
	attemptsLeft int64,
) domain.CreateSessionError {
	err := tx.Commit(ctx)
	if err != nil {
		return unexpected("commit", err)
	}
	return domain.SignInCodeMismatchError{AttemptsLeft: attemptsLeft}
}

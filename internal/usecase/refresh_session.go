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

// RefreshSession trades a key for a fresh pair. The session is not
// replaced: it lasts from signing in until signing out, and these are
// only the keys to it.
//
// Nothing is asked of the caller beyond the key itself, and nothing can
// be: an access token is expired at exactly the moment renewal is
// wanted, so the door this comes through stands open and everything a
// door would have checked is checked here instead — that the key is
// known, unspent, in date, and that the session behind it is still open.
type RefreshSession struct {
	pool *pgxpool.Pool
	seal Seal

	token domain.RefreshToken
}

func NewRefreshSession(
	pool *pgxpool.Pool,
	seal Seal,
	token domain.RefreshToken,
) RefreshSession {
	return RefreshSession{pool: pool, seal: seal, token: token}
}

// Exec spends the key shown and writes its replacement, both inside one
// transaction. A key spent with no replacement behind it would leave its
// holder locked out of a session that is still open, holding something
// that no longer works.
//
// The token is stamped before the commit for the same reason it is when
// signing in: stamping can fail, and failing inside the transaction puts
// the spent key back, so the holder simply tries again.
func (u RefreshSession) Exec(ctx context.Context) (domain.Pass, domain.RefreshSessionError) {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return domain.Pass{}, unexpected("begin transaction", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	spent, err := db.New(tx).SpendRefreshToken(ctx, u.token.Hash())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Pass{}, u.refuse(ctx, tx)
		}
		return domain.Pass{}, unexpected("spend refresh token", err)
	}

	refresh, err := domain.NewRefreshToken()
	if err != nil {
		return domain.Pass{}, unexpected("draw refresh token", err)
	}

	refreshExpiresIn, err := db.New(tx).InsertRefreshToken(ctx, db.InsertRefreshTokenParams{
		SessionID: spent.SessionID,
		Hash:      refresh.Hash(),
	})
	if err != nil {
		return domain.Pass{}, unexpected("insert refresh token", err)
	}

	access, accessExpiresIn, err := u.seal.Stamp(
		domain.ID(spent.UserID),
		domain.ID(spent.SessionID),
	)
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

// refuse turns a key that opened nothing into the answer, and along the
// way finds out whether the failure was the one that has to be acted on.
//
// A key already spent means a copy of it is loose, and the session ends
// for everyone holding one. Because that write has to outlive the
// refusal, this path commits where every other way out of Exec rolls
// back: left to the deferred rollback, the theft would be noticed and
// then quietly forgotten.
//
// Only this path pays for the second question. On the way out the answer
// is needed anyway, and on the way in it is not.
func (u RefreshSession) refuse(ctx context.Context, tx pgx.Tx) domain.RefreshSessionError {
	spent, err := db.New(tx).GetSpentRefreshToken(ctx, u.token.Hash())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Unknown, out of date, or belonging to a session already
			// closed. Nothing to end and nothing to record.
			return domain.RefreshTokenNotUsableError{}
		}
		return unexpected("get spent refresh token", err)
	}

	// The two carry the same pair, so the conversion is exact rather than
	// convenient: Go refuses it the moment either side gains a field or
	// changes one, which is the check a hand-written literal would not
	// have given.
	err = db.New(tx).EndSession(ctx, db.EndSessionParams(spent))
	if err != nil {
		return unexpected("end session", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return unexpected("commit", err)
	}
	return domain.RefreshTokenNotUsableError{}
}

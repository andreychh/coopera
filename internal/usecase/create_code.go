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

// CreateCode sends a one-time code to an address so that whoever reads
// that mailbox can sign in.
//
// It asks nothing about who owns the address, and there is nowhere in it
// to find out: the answer is the same whether an account stands behind
// the address or not, and an account is only born later, when the code
// comes back.
type CreateCode struct {
	pool *pgxpool.Pool
	post Post

	email domain.Email
}

func NewCreateCode(pool *pgxpool.Pool, post Post, email domain.Email) CreateCode {
	return CreateCode{pool: pool, post: post, email: email}
}

// Exec needs no transaction: the limits are checked by the same
// statement that writes the code, so no second request can pass a limit
// that was true a moment ago and false by the time of the write.
//
// The letter goes out after the write, which means a code can exist that
// nobody received — the person is told the sending failed and, having
// spent one of five, waits a minute. The other order is worse: a letter
// no record backs would let its reader in against a code the system does
// not know.
func (u CreateCode) Exec(ctx context.Context) (domain.CodeDelivery, domain.CreateCodeError) {
	code, err := domain.NewSignInCode()
	if err != nil {
		return domain.CodeDelivery{}, unexpected("draw code", err)
	}

	expiresIn, err := db.New(u.pool).InsertSignInCode(ctx, db.InsertSignInCodeParams{
		Email: u.email.String(),
		Code:  code.String(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.CodeDelivery{}, u.tooSoon(ctx)
		}
		return domain.CodeDelivery{}, unexpected("insert sign-in code", err)
	}

	err = u.post.Send(ctx, u.email, code)
	if err != nil {
		return domain.CodeDelivery{}, unexpected("send code", err)
	}

	retryAfter, err := u.retryAfter(ctx)
	if err != nil {
		return domain.CodeDelivery{}, unexpected("read retry after", err)
	}

	return domain.CodeDelivery{
		ExpiresIn:  time.Duration(expiresIn) * time.Second,
		RetryAfter: retryAfter,
	}, nil
}

// tooSoon builds the refusal for an address that has asked too often.
// Only this path pays for the extra question: on the way out the answer
// is needed anyway, and on the way in it is not known.
func (u CreateCode) tooSoon(ctx context.Context) domain.CreateCodeError {
	retryAfter, err := u.retryAfter(ctx)
	if err != nil {
		return unexpected("read retry after", err)
	}
	return domain.TooManyCodesError{RetryAfter: retryAfter}
}

// retryAfter says when the address may ask for a code again. The same
// question serves both outcomes: after a code was sent it is a minute
// away, unless that code was the fifth within the hour, in which case it
// is the oldest one falling out of the window that decides.
func (u CreateCode) retryAfter(ctx context.Context) (time.Duration, error) {
	seconds, err := db.New(u.pool).SignInCodeRetryAfter(ctx, u.email.String())
	if err != nil {
		return 0, err
	}
	return time.Duration(seconds) * time.Second, nil
}

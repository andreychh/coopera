// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/andreychh/coopera/internal/db"
	"github.com/andreychh/coopera/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RevokeInviteLink revokes an invite link. Only the owner of the
// team it belongs to may revoke it.
type RevokeInviteLink struct {
	pool *pgxpool.Pool

	actorID domain.ID
	code    domain.Code
}

func NewRevokeInviteLink(
	pool *pgxpool.Pool,
	actorID domain.ID,
	code domain.Code,
) RevokeInviteLink {
	return RevokeInviteLink{pool: pool, actorID: actorID, code: code}
}

// Exec reads the link, judges it and only then writes, which is why the
// read locks the row: nothing else may revoke it in between and leave
// this call reporting success over someone else's work.
func (u RevokeInviteLink) Exec(ctx context.Context) domain.RevokeInviteLinkError {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return unexpected("begin transaction", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	link, err := db.New(tx).GetInviteLinkByCode(ctx, string(u.code))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.InviteLinkNotFoundError{Code: u.code}
		}
		return unexpected("get invite link", err)
	}

	// Ownership is settled before the link's condition is, so a stranger
	// holding a code learns nothing about it beyond being turned away.
	owner, err := db.New(tx).IsTeamOwner(ctx, db.IsTeamOwnerParams{
		TeamID: link.TeamID,
		UserID: uuid.UUID(u.actorID),
	})
	if err != nil {
		return unexpected("check ownership", err)
	}
	if !owner {
		return domain.NotTeamOwnerError{TeamID: domain.ID(link.TeamID)}
	}

	// Only an already revoked link is a conflict; one that merely lapsed
	// can still be revoked.
	state := domain.NewLinkState(
		(*domain.DateTime)(link.ExpiresAt),
		(*domain.DateTime)(link.RevokedAt),
		domain.DateTime(time.Now()),
	)
	if _, revoked := state.(domain.Revoked); revoked {
		return domain.InviteLinkAlreadyRevokedError{Code: u.code}
	}

	err = db.New(tx).RevokeInviteLink(ctx, link.ID)
	if err != nil {
		return unexpected("revoke invite link", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return unexpected("commit", err)
	}
	return nil
}

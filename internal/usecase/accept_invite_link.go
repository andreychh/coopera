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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AcceptInviteLink adds actor as a member of the team an invite
// link belongs to. Accepting is idempotent: someone who already belongs
// to the team is asking for something that already holds, so the call
// succeeds and changes nothing.
type AcceptInviteLink struct {
	pool *pgxpool.Pool

	actorID domain.ID
	code    domain.Code
}

func NewAcceptInviteLink(
	pool *pgxpool.Pool,
	actorID domain.ID,
	code domain.Code,
) AcceptInviteLink {
	return AcceptInviteLink{pool: pool, actorID: actorID, code: code}
}

// Exec reports the team joined and whether joining actually happened:
// false means actor was already a member, which is a success with
// nothing changed rather than a failure.
//
// Unlike creating a link, this cannot be one statement. Whether a link
// is usable is derived in Go, because listing links has to show that
// derivation too, and restating it in SQL here would make a second copy
// of the one rule. So the link is read, judged and then acted on — and
// the read locks the row, which is what stops a revoke from committing
// in between.
func (u AcceptInviteLink) Exec(
	ctx context.Context,
) (domain.Team, bool, domain.AcceptInviteLinkError) {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return domain.Team{}, false, unexpected("begin transaction", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	link, err := db.New(tx).GetInviteLinkByCode(ctx, string(u.code))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Team{}, false, domain.InviteLinkNotFoundError{Code: u.code}
		}
		return domain.Team{}, false, unexpected("get invite link", err)
	}

	state := domain.NewLinkState(
		(*domain.DateTime)(link.ExpiresAt),
		(*domain.DateTime)(link.RevokedAt),
		domain.DateTime(time.Now()),
	)
	if _, active := state.(domain.Active); !active {
		return domain.Team{}, false, domain.InviteLinkNotUsableError{Code: u.code}
	}

	joined, joinErr := u.join(ctx, tx, link.TeamID)
	if joinErr != nil {
		return domain.Team{}, false, joinErr
	}

	// The count is of people who joined, so opening a link again does not
	// raise it.
	if joined {
		err = db.New(tx).IncrementInviteLinkUseCount(ctx, link.ID)
		if err != nil {
			return domain.Team{}, false, unexpected("increment use count", err)
		}
	}

	team, err := db.New(tx).GetTeam(ctx, link.TeamID)
	if err != nil {
		return domain.Team{}, false, unexpected("get team", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.Team{}, false, unexpected("commit", err)
	}

	return domain.Team{
		ID:        domain.ID(team.ID),
		Name:      domain.TeamName(team.Name),
		CreatedAt: domain.DateTime(team.CreatedAt),
	}, joined, nil
}

// join makes actor a member of teamID and reports whether that changed
// anything. No row back means they were already an active member: the
// upsert skips their row on purpose, and that absence is the idempotent
// outcome rather than a failure.
func (u AcceptInviteLink) join(
	ctx context.Context,
	tx pgx.Tx,
	teamID uuid.UUID,
) (bool, domain.AcceptInviteLinkError) {
	_, err := db.New(tx).JoinTeam(ctx, db.JoinTeamParams{
		TeamID: teamID,
		UserID: uuid.UUID(u.actorID),
	})
	if err == nil {
		return true, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	pgErr, isPg := errors.AsType[*pgconn.PgError](err)
	if isPg && pgErr.Code == domain.ForeignKeyViolation &&
		pgErr.ConstraintName == domain.MembersUserIDForeignKey {
		return false, domain.UserNotFoundError{ID: u.actorID}
	}
	return false, unexpected("join team", err)
}

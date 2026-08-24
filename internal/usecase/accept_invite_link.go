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
// link belongs to. Joining happens once: someone who already belongs to
// the team is refused with [domain.AlreadyMemberError] rather than
// answered as though they had just joined.
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

// Exec reports the team joined. Being a member already is not one of the
// ways this succeeds: it comes back as [domain.AlreadyMemberError],
// carrying the team so that the refusal still answers what was asked.
//
// Unlike creating a link, this cannot be one statement. Whether a link
// is usable is derived in Go, because listing links has to show that
// derivation too, and restating it in SQL here would make a second copy
// of the one rule. So the link is read, judged and then acted on — and
// the read locks the row, which is what stops a revoke from committing
// in between.
func (u AcceptInviteLink) Exec(
	ctx context.Context,
) (domain.Team, domain.AcceptInviteLinkError) {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return domain.Team{}, unexpected("begin transaction", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	link, err := db.New(tx).GetInviteLinkByCode(ctx, string(u.code))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Team{}, domain.InviteLinkNotFoundError{Code: u.code}
		}
		return domain.Team{}, unexpected("get invite link", err)
	}

	state := domain.NewLinkState(
		(*domain.DateTime)(link.ExpiresAt),
		(*domain.DateTime)(link.RevokedAt),
		domain.DateTime(time.Now()),
	)
	if _, active := state.(domain.Active); !active {
		return domain.Team{}, domain.InviteLinkNotUsableError{Code: u.code}
	}

	joinErr := u.join(ctx, tx, link.TeamID)
	if joinErr != nil {
		return domain.Team{}, joinErr
	}

	// Reached only by someone who has just joined, so the count of people
	// who joined always moves with it.
	err = db.New(tx).IncrementInviteLinkUseCount(ctx, link.ID)
	if err != nil {
		return domain.Team{}, unexpected("increment use count", err)
	}

	team, err := db.New(tx).GetTeam(ctx, link.TeamID)
	if err != nil {
		return domain.Team{}, unexpected("get team", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.Team{}, unexpected("commit", err)
	}

	return domain.Team{
		ID:        domain.ID(team.ID),
		Name:      domain.TeamName(team.Name),
		CreatedAt: domain.DateTime(team.CreatedAt),
	}, nil
}

// join makes actor a member of teamID. No row back means they were
// already an active member: the upsert skips their row on purpose, and
// that absence is what tells the two apart — there is nothing to add, so
// the asking is refused rather than granted twice.
func (u AcceptInviteLink) join(
	ctx context.Context,
	tx pgx.Tx,
	teamID uuid.UUID,
) domain.AcceptInviteLinkError {
	_, err := db.New(tx).JoinTeam(ctx, db.JoinTeamParams{
		TeamID: teamID,
		UserID: uuid.UUID(u.actorID),
	})
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.AlreadyMemberError{TeamID: domain.ID(teamID)}
	}

	pgErr, isPg := errors.AsType[*pgconn.PgError](err)
	if isPg && pgErr.Code == domain.ForeignKeyViolation &&
		pgErr.ConstraintName == domain.MembersUserIDForeignKey {
		return domain.UserNotFoundError{ID: u.actorID}
	}
	return unexpected("join team", err)
}

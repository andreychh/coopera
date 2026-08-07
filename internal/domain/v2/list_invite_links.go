// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package v2

import (
	"context"
	"fmt"
	"time"

	"github.com/andreychh/coopera/internal/db"
	"github.com/andreychh/coopera/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ListInviteLinksUsecase lists a team's invite links, optionally filtered
// by status. Only the team's owner may list them.
type ListInviteLinksUsecase struct {
	pool *pgxpool.Pool

	actorID domain.ID
	teamID  domain.ID
	status  *domain.LinkStatus
}

func NewListInviteLinksUsecase(
	pool *pgxpool.Pool,

	actorID domain.ID,
	teamID domain.ID,
	status *domain.LinkStatus,
) ListInviteLinksUsecase {
	return ListInviteLinksUsecase{
		pool:    pool,
		actorID: actorID,
		teamID:  teamID,
		status:  status,
	}
}

// Exec asks about ownership separately from reading the links, because
// an empty result cannot carry both answers: an owner whose team has no
// links deserves an empty list, not the refusal a stranger gets.
//
// The two reads are not in a transaction. Losing ownership between them
// would only mean showing links the actor was entitled to a moment
// earlier, which is what any read races with anyway.
func (u ListInviteLinksUsecase) Exec(
	ctx context.Context,
) ([]domain.InviteLink, domain.ListInviteLinksError) {
	owner, err := db.New(u.pool).IsTeamOwner(ctx, db.IsTeamOwnerParams{
		TeamID: uuid.UUID(u.teamID),
		UserID: uuid.UUID(u.actorID),
	})
	if err != nil {
		return nil, domain.UnexpectedError{Err: fmt.Errorf("check ownership: %w", err)}
	}
	if !owner {
		return nil, domain.NotTeamOwnerError{TeamID: u.teamID}
	}

	rows, err := db.New(u.pool).ListInviteLinksByTeam(ctx, uuid.UUID(u.teamID))
	if err != nil {
		return nil, domain.UnexpectedError{Err: fmt.Errorf("list invite links: %w", err)}
	}

	// One instant for the whole list, so two links with the same deadline
	// can't come back in different states.
	now := domain.DateTime(time.Now())

	links := make([]domain.InviteLink, 0, len(rows))
	for _, row := range rows {
		state := domain.NewLinkState(
			(*domain.DateTime)(row.ExpiresAt),
			(*domain.DateTime)(row.RevokedAt),
			now,
		)
		if u.status != nil && !u.status.Matches(state) {
			continue
		}

		links = append(links, domain.InviteLink{
			Code:      domain.Code(row.Code),
			UseCount:  row.UseCount,
			State:     state,
			CreatedAt: domain.DateTime(row.CreatedAt),
		})
	}
	return links, nil
}

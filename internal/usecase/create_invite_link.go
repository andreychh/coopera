// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andreychh/coopera/internal/db"
	"github.com/andreychh/coopera/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateInviteLink creates a link that can be used to join a team.
// Only the team's owner may create one.
type CreateInviteLink struct {
	pool *pgxpool.Pool

	actorID  domain.ID
	teamID   domain.ID
	validity domain.Validity
}

func NewCreateInviteLink(
	pool *pgxpool.Pool,
	actorID domain.ID,
	teamID domain.ID,
	validity domain.Validity,
) CreateInviteLink {
	return CreateInviteLink{
		pool:     pool,
		actorID:  actorID,
		teamID:   teamID,
		validity: validity,
	}
}

func (u CreateInviteLink) Exec(
	ctx context.Context,
) (domain.InviteLink, domain.CreateInviteLinkError) {
	code, err := domain.NewCode()
	if err != nil {
		return domain.InviteLink{}, domain.UnexpectedError{Err: err}
	}

	var expiresAt *domain.DateTime
	switch v := u.validity.(type) {
	case domain.Until:
		expiresAt = new(v.Time)
	case domain.Forever:
		expiresAt = nil
	default:
		return domain.InviteLink{}, domain.UnexpectedError{
			Err: fmt.Errorf("unknown validity: %T", u.validity),
		}
	}

	inserted, err := db.New(u.pool).InsertInviteLinkAsOwner(ctx, db.InsertInviteLinkAsOwnerParams{
		UserID:    uuid.UUID(u.actorID),
		TeamID:    uuid.UUID(u.teamID),
		Code:      string(code),
		ExpiresAt: (*time.Time)(expiresAt),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.InviteLink{}, domain.NotTeamOwnerError{TeamID: u.teamID}
		}
		return domain.InviteLink{}, domain.UnexpectedError{
			Err: fmt.Errorf("insert invite link: %w", err),
		}
	}

	return domain.InviteLink{
		Code:     domain.Code(inserted.Code),
		UseCount: inserted.UseCount,
		State: domain.NewLinkState(
			(*domain.DateTime)(inserted.ExpiresAt),
			(*domain.DateTime)(inserted.RevokedAt),
			domain.DateTime(time.Now()),
		),
		CreatedAt: domain.DateTime(inserted.CreatedAt),
	}, nil
}

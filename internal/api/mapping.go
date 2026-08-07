// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"fmt"

	"github.com/andreychh/coopera/internal/domain"
)

// This file holds pure domain-to-API mapping functions: they only shape
// data already fetched by the caller, never do I/O themselves (no ctx,
// no db/pool), and may call each other.

func newTeams(infos []domain.Team) []Team {
	teams := make([]Team, 0, len(infos))
	for _, info := range infos {
		teams = append(teams, newTeam(info))
	}
	return teams
}

func newTeam(info domain.Team) Team {
	return Team{
		Id:        info.ID.String(),
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt.String(),
	}
}

func newInviteLinks(links []domain.InviteLink) ([]InviteLink, error) {
	items := make([]InviteLink, 0, len(links))
	for _, link := range links {
		item, err := newInviteLink(link)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func newInviteLink(link domain.InviteLink) (InviteLink, error) {
	state, err := newInviteLinkState(link.State)
	if err != nil {
		return InviteLink{}, err
	}
	return InviteLink{
		Code:      link.Code.String(),
		CreatedAt: link.CreatedAt.String(),
		State:     state,
		UseCount:  int(link.UseCount),
	}, nil
}

// newInviteLinkState renders a [domain.LinkState] as the discriminated
// union the spec declares. Each case carries exactly the field its
// schema requires, so there is no longer a way to answer with an empty
// expired_at or revoked_at.
func newInviteLinkState(state domain.LinkState) (InviteLinkState, error) {
	var out InviteLinkState
	switch s := state.(type) {
	case domain.Active:
		active := ActiveInviteLinkState{
			Status:    ActiveInviteLinkStateStatusActive,
			ExpiresAt: nil,
		}
		switch v := s.Validity.(type) {
		case domain.Until:
			active.ExpiresAt = new(v.Time.String())
		case domain.Forever:
			// Reported as an explicit null: the link has no deadline.
		default:
			return out, fmt.Errorf("unknown validity: %T", s.Validity)
		}
		err := out.FromActiveInviteLinkState(active)
		return out, err

	case domain.Expired:
		err := out.FromExpiredInviteLinkState(ExpiredInviteLinkState{
			Status:    Expired,
			ExpiredAt: s.At.String(),
		})
		return out, err

	case domain.Revoked:
		err := out.FromRevokedInviteLinkState(RevokedInviteLinkState{
			Status:    Revoked,
			RevokedAt: s.At.String(),
		})
		return out, err
	}
	return out, fmt.Errorf("unknown link state: %T", state)
}

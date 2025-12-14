package http

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpMembers struct {
	teamID int64
	client transport.Client
}

func (h httpMembers) All(ctx context.Context) ([]domain.Member, error) {
	resp := findTeamResponse{}
	err := h.client.Get(
		ctx,
		transport.URL("teams").With("team_id", strconv.FormatInt(h.teamID, 10)).String(),
		&resp,
	)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", h.teamID, err)
	}
	members := make([]domain.Member, 0, len(resp.Members))
	for _, member := range resp.Members {
		members = append(
			members,
			Member(member.MemberId, member.Username, domain.MemberRole(member.Role), h.teamID, h.client),
		)
	}
	return members, nil
}

func (h httpMembers) MemberWithUsername(ctx context.Context, username string) (domain.Member, bool, error) {
	resp := findTeamResponse{}
	err := h.client.Get(
		ctx,
		transport.URL("teams").With("team_id", strconv.FormatInt(h.teamID, 10)).String(),
		&resp,
	)
	if err != nil {
		return nil, false, fmt.Errorf("getting team %d: %w", h.teamID, err)
	}
	for _, member := range resp.Members {
		if member.Username == username {
			return Member(member.MemberId, member.Username, domain.MemberRole(member.Role), h.teamID, h.client), true, nil
		}
	}
	return nil, false, nil
}

func Members(teamID int64, client transport.Client) domain.Members {
	return httpMembers{
		teamID: teamID,
		client: client,
	}
}

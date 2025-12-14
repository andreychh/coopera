package http

import (
	"context"
	"encoding/json"
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
	data, err := h.client.Get(
		ctx,
		transport.NewOutcomingURL("teams").
			With("team_id", strconv.FormatInt(h.teamID, 10)).
			String(),
	)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", h.teamID, err)
	}
	resp := struct {
		Members []struct {
			ID       int64             `json:"member_id"`
			Username string            `json:"username"`
			UserID   int64             `json:"user_id"`
			Role     domain.MemberRole `json:"role"`
		} `json:"members"`
	}{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	members := make([]domain.Member, 0, len(resp.Members))
	for _, m := range resp.Members {
		members = append(members, Member(m.ID, m.UserID, h.teamID, m.Username, m.Role, h.client))
	}
	return members, nil
}

func Members(teamID int64, client transport.Client) domain.Members {
	return httpMembers{
		teamID: teamID,
		client: client,
	}
}

package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpTeam struct {
	id     int64
	name   string
	client transport.Client
}

func (h httpTeam) ID() int64 {
	return h.id
}

func (h httpTeam) Name() string {
	return h.name
}

func (h httpTeam) Stats(ctx context.Context) (domain.TeamStats, error) {
	panic("not implemented")
}

func (h httpTeam) AddMember(ctx context.Context, userID int64) (domain.Member, error) {
	req := createMemberRequest{
		TeamId: h.id,
		UserId: userID,
	}
	resp := createMemberResponse{}
	err := h.client.Post(ctx, "memberships", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("adding member to team %d: %w", h.id, err)
	}
	return nil, nil
	// username, role, err := h.member(ctx, resp.Id)
	// if err != nil {
	// 	return nil, fmt.Errorf("getting username for member %d in team %d: %w", userID, h.id, err)
	// }
	// return Member(resp.Id, username, role, h.id, h.client), nil
}

func (h httpTeam) Members(ctx context.Context) (domain.Members, error) {
	return Members(h.id, h.client), nil
}

func (h httpTeam) Tasks(ctx context.Context) (domain.Tasks, error) {
	return TeamTasks(h.id, h.client), nil
}

func (h httpTeam) member(ctx context.Context, memberID int64) (string, domain.MemberRole, error) {
	resp := findTeamResponse{}
	err := h.client.Get(
		ctx,
		transport.URL("teams").With("team_id", strconv.FormatInt(h.id, 10)).String(),
		&resp,
	)
	var apiErr transport.APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
		return "", "", fmt.Errorf("team %d not found", h.id)
	}
	if err != nil {
		return "", "", fmt.Errorf("getting team %d: %w", h.id, err)
	}
	for _, member := range resp.Members {
		if member.MemberId == memberID {
			return member.Username, domain.MemberRole(member.Role), nil
		}
	}
	return "", "", fmt.Errorf("member %d not found in team %d", memberID, h.id)
}

func Team(id int64, name string, client transport.Client) domain.Team {
	return httpTeam{
		id:     id,
		name:   name,
		client: client,
	}
}

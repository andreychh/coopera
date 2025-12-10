package http

import (
	"context"
	jsn "encoding/json"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/memory"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
	"github.com/tidwall/gjson"
)

type httpTeam struct {
	id     int64
	client transport.Client
}

func (h httpTeam) Exists(ctx context.Context) (bool, error) {
	panic("not implemented")
}

func (h httpTeam) Details(ctx context.Context) (domain.TeamDetails, error) {
	data, err := h.client.Get(ctx, transport.NewOutcomingURL("teams").
		With("team_id", strconv.FormatInt(h.id, 10)).
		String(),
	)
	if err != nil {
		return nil, err
	}
	return memory.TeamDetails(
		gjson.GetBytes(data, "id").Int(),
		gjson.GetBytes(data, "name").String(),
	), nil
}

func (h httpTeam) AddMember(ctx context.Context, user domain.User) (domain.Member, error) {
	details, err := user.Details(ctx)
	if err != nil {
		return nil, err
	}
	payload, err := json.Object(json.Fields{
		"team_id": json.I64(h.id),
		"user_id": json.I64(details.ID()),
	}).Marshal()
	if err != nil {
		return nil, fmt.Errorf("marshaling payload: %w", err)
	}
	_, err = h.client.Post(ctx, "memberships", payload)
	if err != nil {
		return nil, err
	}
	return httpMember{
		userID:     details.ID(),
		teamID:     h.id,
		dataSource: h.client,
	}, nil
}

func (h httpTeam) Members(ctx context.Context) ([]domain.Member, error) {
	data, err := h.client.Get(ctx, transport.NewOutcomingURL("teams").
		With("team_id", strconv.FormatInt(h.id, 10)).
		String())
	if err != nil {
		return nil, err
	}
	membersJSON := gjson.GetBytes(data, "members").Raw
	if membersJSON == "" || membersJSON == "null" {
		return nil, fmt.Errorf("field 'members' not found in response")
	}
	var membersSlice []struct {
		ID   int64  `json:"member_id"`
		Role string `json:"role"`
	}
	err = jsn.Unmarshal([]byte(membersJSON), &membersSlice)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling members: %w", err)
	}
	members := make([]domain.Member, 0, len(membersSlice))
	for _, m := range membersSlice {
		members = append(members, cachedMember{
			userID: m.ID,
			role:   m.Role,
			client: h.client,
		})
	}
	return members, nil
}

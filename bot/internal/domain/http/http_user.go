package http

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/transport"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
	"github.com/tidwall/gjson"
)

type httpUser struct {
	id     int64
	client transport.Client
}

func (h httpUser) CreatedTeams(ctx context.Context) ([]domain.Team, error) {
	panic("not implemented")
}

func (h httpUser) CreateTeam(ctx context.Context, name string) (domain.Team, error) {
	payload, err := json.Object(json.Fields{"user_id": json.I64(h.id), "name": json.Str(name)}).Marshal()
	if err != nil {
		return nil, fmt.Errorf("marshaling payload: %w", err)
	}
	data, err := h.client.Post(ctx, "teams", payload)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}
	teamID := gjson.GetBytes(data, "id")
	if !teamID.Exists() {
		return nil, fmt.Errorf("field 'id' not found in response")
	}
	return nil, nil
}

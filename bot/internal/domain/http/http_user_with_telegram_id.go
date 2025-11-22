package http

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"

	jsn "encoding/json"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/transport"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
	"github.com/tidwall/gjson"
)

type httpUserWithTelegramID struct {
	telegramID int64
	client     transport.Client
	id         atomic.Int64
}

func (h *httpUserWithTelegramID) CreatedTeams(ctx context.Context) ([]domain.Team, error) {
	data, err := h.client.Get(ctx, transport.NewOutcomingURL("users").
		With("telegram_id", strconv.FormatInt(h.telegramID, 10)).
		String(),
	)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}
	id := gjson.GetBytes(data, "id")
	if !id.Exists() {
		return nil, fmt.Errorf("field 'id' not found in response")
	}
	h.id.Store(id.Int())
	teamsJSON := gjson.GetBytes(data, "teams").Raw
	if teamsJSON == "" || teamsJSON == "null" {
		return nil, fmt.Errorf("field 'teams' not found in response")
	}
	var teamsSlice []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	err = jsn.Unmarshal([]byte(teamsJSON), &teamsSlice)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling teams: %w", err)
	}
	teams := make([]domain.Team, 0, len(teamsSlice))
	for _, t := range teamsSlice {
		teams = append(teams, cachedTeam{
			id:     t.ID,
			name:   t.Name,
			client: h.client,
		})
	}
	return teams, nil
}

func (h *httpUserWithTelegramID) CreateTeam(ctx context.Context, name string) (domain.Team, error) {
	id, err := h.resolveID(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving user ID: %w", err)
	}
	payload, err := json.Object(json.Fields{"user_id": json.I64(id), "name": json.Str(name)}).Marshal()
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

func (h *httpUserWithTelegramID) resolveID(ctx context.Context) (int64, error) {
	if id := h.id.Load(); id != 0 {
		return id, nil
	}
	data, err := h.client.Get(ctx, transport.NewOutcomingURL("users").
		With("telegram_id", strconv.FormatInt(h.telegramID, 10)).
		String(),
	)
	if err != nil {
		return 0, fmt.Errorf("getting user: %w", err)
	}
	id := gjson.GetBytes(data, "id")
	if !id.Exists() {
		return 0, fmt.Errorf("field 'id' not found in response")
	}
	h.id.Store(id.Int())
	return id.Int(), nil
}

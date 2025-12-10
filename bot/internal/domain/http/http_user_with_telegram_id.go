package http

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/memory"
	domaintransport "github.com/andreychh/coopera-bot/internal/domain/transport"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
	"github.com/tidwall/gjson"
)

type httpUserWithTelegramID struct {
	telegramID int64
	client     domaintransport.Client
	id         atomic.Int64
}

func (h *httpUserWithTelegramID) Details(ctx context.Context) (domain.UserDetails, error) {
	id, err := h.resolveID(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving user ID: %w", err)
	}
	return memory.UserDetails(id), nil
}

func (h *httpUserWithTelegramID) CreatedTeams() domain.Teams {
	panic("not implemented")
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
	data, err := h.client.Get(ctx, domaintransport.NewOutcomingURL("users").
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

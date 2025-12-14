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

type httpCommunity struct {
	client transport.Client
}

func (h httpCommunity) CreateUser(ctx context.Context, tgID int64, tgUsername string) (domain.User, error) {
	req := createUserRequest{
		TelegramId: tgID,
		Username:   tgUsername,
	}
	resp := createUserResponse{}
	err := h.client.Post(ctx, "users", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	return User(resp.Id, resp.Username, h.client), nil
}

func (h httpCommunity) UserWithID(ctx context.Context, id int64) (domain.User, bool, error) {
	return h.user(ctx, "id", strconv.FormatInt(id, 10))
}

func (h httpCommunity) UserWithTelegramID(ctx context.Context, tgID int64) (domain.User, bool, error) {
	return h.user(ctx, "telegram_id", strconv.FormatInt(tgID, 10))
}

func (h httpCommunity) UserWithUsername(ctx context.Context, tgUsername string) (domain.User, bool, error) {
	return h.user(ctx, "username", tgUsername)
}

func (h httpCommunity) user(ctx context.Context, key, value string) (domain.User, bool, error) {
	resp := findUserResponse{}
	err := h.client.Get(
		ctx,
		transport.URL("users").With(key, value).String(),
		&resp,
	)
	var apiErr transport.APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("getting user by %s=%s: %w", key, value, err)
	}
	return User(resp.Id, resp.Username, h.client), true, nil
}

func (h httpCommunity) Team(ctx context.Context, id int64) (domain.Team, bool, error) {
	resp := findTeamResponse{}
	err := h.client.Get(
		ctx,
		transport.URL("teams").With("team_id", strconv.FormatInt(id, 10)).String(),
		&resp,
	)
	var apiErr transport.APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("getting team %d: %w", id, err)
	}
	return Team(resp.Id, resp.Name, h.client), true, nil
}

func (h httpCommunity) Task(ctx context.Context, id int64) (domain.Task, bool, error) {
	var resp []findTasksResponse
	err := h.client.Get(
		ctx,
		transport.URL("tasks").With("task_id", strconv.FormatInt(id, 10)).String(),
		&resp,
	)
	var apiErr transport.APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("getting task %d: %w", id, err)
	}
	return RespToTask(resp[0], h.client), true, nil
}

func Community(client transport.Client) domain.Community {
	return httpCommunity{
		client: client,
	}
}

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpTask struct {
	id     int64
	title  string
	points int
	teamID int64
	status string
	client transport.Client
}

func (h httpTask) ID() int64 {
	return h.id
}

func (h httpTask) Title() string {
	return h.title
}

func (h httpTask) Points() int {
	return h.points
}

func (h httpTask) Status() string {
	return h.status
}

func (h httpTask) Team(ctx context.Context) (domain.Team, error) {
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
		Name string `json:"name"`
	}{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	return Team(h.teamID, resp.Name, h.client), nil
}

func Task(id int64, title string, points int, status string, teamID int64, client transport.Client) domain.Task {
	return httpTask{
		id:     id,
		title:  title,
		points: points,
		teamID: teamID,
		status: status,
		client: client,
	}
}

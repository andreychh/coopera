package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpMemberTasks struct {
	memberID int64
	teamID   int64
	client   transport.Client
}

func (h httpMemberTasks) All(ctx context.Context) ([]domain.Task, error) {
	data, err := h.client.Get(ctx, transport.NewOutcomingURL("tasks").
		With("team_id", strconv.FormatInt(h.teamID, 10)).
		String(),
	)
	if err != nil {
		return nil, fmt.Errorf("getting tasks for team %d: %w", h.teamID, err)
	}
	var resp []task
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling data: %w", err)
	}
	tasks := make([]domain.Task, 0, len(resp))
	for _, t := range resp {
		if t.AssignedToMember == h.memberID {
			tasks = append(tasks, Task(t.ID, t.Title, t.Points, t.Status, t.TeamID, h.client))
		}
	}
	return tasks, nil
}

func (h httpMemberTasks) Empty(ctx context.Context) (bool, error) {
	data, err := h.client.Get(ctx, transport.NewOutcomingURL("tasks").
		With("team_id", strconv.FormatInt(h.teamID, 10)).
		String(),
	)
	if err != nil {
		return false, fmt.Errorf("getting tasks for team %d: %w", h.teamID, err)
	}
	var resp []task
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return false, fmt.Errorf("unmarshaling data: %w", err)
	}
	for _, t := range resp {
		if t.AssignedToMember == h.memberID {
			return false, nil
		}
	}
	return true, nil
}

func MemberTasks(memberID int64, teamID int64, client transport.Client) domain.Tasks {
	return httpMemberTasks{
		memberID: memberID,
		teamID:   teamID,
		client:   client,
	}
}

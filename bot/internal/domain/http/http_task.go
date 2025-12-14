package http

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpTask struct {
	id               int64
	title            string
	description      string
	points           *int
	createdAt        time.Time
	createdByMember  int64
	assignedToMember *int64
	teamID           int64
	status           domain.TaskStatus
	client           transport.Client
}

func (h httpTask) ID() int64 {
	return h.id
}

func (h httpTask) Title() string {
	return h.title
}

func (h httpTask) Description() string {
	return h.description
}

func (h httpTask) Points() (int, bool) {
	if h.points == nil {
		return 0, false
	}
	return *h.points, true
}

func (h httpTask) Status() domain.TaskStatus {
	return h.status
}

func (h httpTask) CreatedAt() time.Time {
	return h.createdAt
}

func (h httpTask) Assignee(ctx context.Context) (domain.Member, bool, error) {
	if h.assignedToMember == nil {
		return nil, false, nil
	}
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
		if member.MemberId == *h.assignedToMember {
			return Member(member.MemberId, member.Username, domain.MemberRole(member.Role), h.teamID, h.client), true, nil
		}
	}
	return nil, false, fmt.Errorf("member %d not found in team %d", *h.assignedToMember, h.teamID)
}

func (h httpTask) Team(ctx context.Context) (domain.Team, error) {
	team, exists, err := Community(h.client).Team(ctx, h.teamID)
	if !exists {
		return nil, fmt.Errorf("team %d not found", h.teamID)
	}
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", h.teamID, err)
	}
	return team, nil
}

func Task(
	id int64,
	title string,
	description string,
	points *int,
	createdAt time.Time,
	createdByMember int64,
	assignedToMember *int64,
	teamID int64,
	status domain.TaskStatus,
	client transport.Client,
) domain.Task {
	return httpTask{
		id:               id,
		title:            title,
		description:      description,
		points:           points,
		createdAt:        createdAt,
		createdByMember:  createdByMember,
		assignedToMember: assignedToMember,
		teamID:           teamID,
		status:           status,
		client:           client,
	}
}

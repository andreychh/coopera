package http

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type httpUser struct {
	id       int64
	username string
	client   transport.Client
}

func (h httpUser) ID() int64 {
	return h.id
}

func (h httpUser) Username() string {
	return h.username
}

func (h httpUser) Stats(ctx context.Context) (domain.UserStats, error) {
	userResp := findUserResponse{}
	err := h.client.Get(
		ctx,
		transport.URL("users").With("id", strconv.FormatInt(h.id, 10)).String(),
		&userResp,
	)
	if err != nil {
		return domain.UserStats{}, fmt.Errorf("getting user %s: %w", h.username, err)
	}
	stats := domain.UserStats{
		Teams: make(map[string]domain.UserTeamStats),
	}
	for _, team := range userResp.Teams {
		memberID, err := h.fetchMemberID(ctx, team.Id)
		if err != nil {
			return domain.UserStats{}, fmt.Errorf("fetching member ID for team %d: %w", team.Id, err)
		}
		var tasksResp []findTasksResponse
		err = h.client.Get(
			ctx,
			transport.URL("tasks").With("member_id", strconv.FormatInt(memberID, 10)).String(),
			&tasksResp,
		)
		if err != nil {
			return domain.UserStats{}, fmt.Errorf("getting tasks for team %d: %w", team.Id, err)
		}
		teamStat := domain.UserTeamStats{}
		for _, t := range tasksResp {
			if t.AssignedToMember == nil || *t.AssignedToMember != memberID {
				continue
			}
			if t.Points == nil {
				continue
			}
			points := *t.Points
			if t.Status == "completed" {
				teamStat.LifetimeContribution.TasksCompleted++
				teamStat.LifetimeContribution.PointsEarned += points
			} else {
				teamStat.ActiveLoad.TasksCount++
				teamStat.ActiveLoad.TotalPoints += points
			}
		}
		stats.Teams[team.Name] = teamStat
	}
	return stats, nil
}

func (h httpUser) fetchMemberID(ctx context.Context, teamID int64) (int64, error) {
	var resp findTeamResponse
	err := h.client.Get(
		ctx,
		transport.URL("teams").With("team_id", strconv.FormatInt(teamID, 10)).String(),
		&resp,
	)
	if err != nil {
		return 0, fmt.Errorf("getting team %d: %w", teamID, err)
	}
	for _, member := range resp.Members {
		if member.Username == h.username {
			return member.MemberId, nil
		}
	}
	return 0, fmt.Errorf("member with username %s not found in team %d", h.username, teamID)
}

func (h httpUser) CreateTeam(ctx context.Context, name string) (domain.Team, error) {
	req := createTeamRequest{
		UserId: h.id,
		Name:   name,
	}
	resp := createTeamResponse{}
	err := h.client.Post(ctx, "teams", req, &resp)
	if err != nil {
		return nil, err
	}
	return Team(resp.Id, resp.Name, h.client), nil
}

func (h httpUser) Teams(ctx context.Context) (domain.Teams, error) {
	return Teams(h.id, h.client), nil
}

func (h httpUser) AssignedTasks(ctx context.Context) (domain.Tasks, error) {
	return UserTasks(h.id, h.username, h.client), nil
}

func User(id int64, username string, client transport.Client) domain.User {
	return httpUser{
		id:       id,
		username: username,
		client:   client,
	}
}

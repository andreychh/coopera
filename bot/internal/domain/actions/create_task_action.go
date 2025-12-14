package actions

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type createTaskAction struct {
	community domain.Community
	forms     forms.Forms
}

func (c createTaskAction) Perform(ctx context.Context, update telegram.Update) error {
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return fmt.Errorf("chat ID not found")
	}
	f := c.forms.Form(chatID)
	// Team ID
	teamIDStr, err := f.Field("team_id").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting team_id: %w", err)
	}
	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("parsing team_id: %w", err)
	}
	// Title
	title, err := f.Field("task_title").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting task_title: %w", err)
	}
	// Description
	description, err := f.Field("task_description").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting task_description: %w", err)
	}
	// Points
	pointsStr, err := f.Field("task_points").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting task_points: %w", err)
	}
	points, err := strconv.Atoi(pointsStr)
	if err != nil {
		return fmt.Errorf("parsing task_points: %w", err)
	}
	// Assignee
	assigneeUsername, err := f.Field("task_assignee").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting task_points: %w", err)
	}
	team, err := c.community.Team(ctx, teamID)
	if err != nil {
		return fmt.Errorf("getting team %d: %w", teamID, err)
	}
	var assignee domain.Member
	if assigneeUsername == "" {
		assignee = domain.NullMember()
	} else {
		user, err := c.community.UserWithUsername(ctx, assigneeUsername)
		if err != nil {
			return fmt.Errorf("getting user with username %q: %w", assigneeUsername, err)
		}
		assignee, err = team.MemberWithUserID(ctx, user.ID())
		if err != nil {
			return fmt.Errorf("getting assignee %q: %w", assigneeUsername, err)
		}
	}
	creator, err := c.community.UserWithTelegramID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("getting creator user: %w", err)
	}
	member, err := team.MemberWithUserID(ctx, creator.ID())
	if err != nil {
		return fmt.Errorf("getting member for user %d in team %d: %w", creator.ID(), teamID, err)
	}
	_, err = member.CreateTask(ctx, title, description, points, assignee)
	return err
}

func CreateTask(community domain.Community, forms forms.Forms) core.Action {
	return createTaskAction{
		community: community,
		forms:     forms,
	}
}

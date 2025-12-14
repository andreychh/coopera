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
	team, exists, err := c.community.Team(ctx, teamID)
	if err != nil {
		return fmt.Errorf("getting team %d: %w", teamID, err)
	}
	if !exists {
		return fmt.Errorf("team with ID %d does not exist", teamID)
	}
	members, err := team.Members(ctx)
	if err != nil {
		return fmt.Errorf("getting members of team %d: %w", teamID, err)
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
	creatorUser, exists, err := c.community.UserWithTelegramID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("getting user with telegram ID %d: %w", chatID, err)
	}
	if !exists {
		return fmt.Errorf("user with telegram ID %d does not exist", chatID)
	}
	creatorMember, exists, err := members.MemberWithUsername(ctx, creatorUser.Username())
	if err != nil {
		return fmt.Errorf("getting member with username %q in team %d: %w", creatorUser.Username(), teamID, err)
	}
	if !exists {
		return fmt.Errorf("user with username %q is not a member of team %d", creatorUser.Username(), teamID)
	}
	if assigneeUsername == "" {
		_, err = creatorMember.CreateUnassigned(ctx, title, description, points)
		if err != nil {
			return fmt.Errorf("creating unassigned task in team %d: %w", teamID, err)
		}
		return nil
	}
	assigneeUser, exists, err := c.community.UserWithUsername(ctx, assigneeUsername)
	if err != nil {
		return fmt.Errorf("getting user with username %q: %w", assigneeUsername, err)
	}
	if !exists {
		return fmt.Errorf("user with username %q does not exist", assigneeUsername)
	}
	_, err = creatorMember.CreateAssigned(ctx, title, description, points, assigneeUser.ID())
	if err != nil {
		return fmt.Errorf("creating assigned task in team %d: %w", teamID, err)
	}
	return nil
}

func CreateTask(community domain.Community, forms forms.Forms) core.Action {
	return createTaskAction{
		community: community,
		forms:     forms,
	}
}

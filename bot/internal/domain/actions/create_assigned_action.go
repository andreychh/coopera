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

type createAssignedAction struct {
	community domain.Community
	forms     forms.Forms
}

func (c createAssignedAction) Perform(ctx context.Context, update telegram.Update) error {
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return fmt.Errorf("chat ID not found")
	}
	form := c.forms.Form(chatID)
	// Team ID
	teamIDStr, err := form.Field("team_id").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting team_id: %w", err)
	}
	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("parsing team_id: %w", err)
	}
	// Title
	title, err := form.Field("task_title").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting task_title: %w", err)
	}
	// Description
	description, err := form.Field("task_description").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting task_description: %w", err)
	}
	// Points
	pointsStr, err := form.Field("task_points").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting task_points: %w", err)
	}
	points, err := strconv.Atoi(pointsStr)
	if err != nil {
		return fmt.Errorf("parsing task_points: %w", err)
	}
	// Assignee
	assigneeUsername, err := form.Field("task_assignee").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting task_assignee: %w", err)
	}
	// Create Task
	creatorUser, exists, err := c.community.UserWithTelegramID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("getting user with telegram ID %d: %w", chatID, err)
	}
	if !exists {
		return fmt.Errorf("user with telegram ID %d does not exist", chatID)
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
	assigneeMember, exists, err := members.MemberWithUsername(ctx, assigneeUsername)
	if err != nil {
		return fmt.Errorf("getting member with username %q in team %d: %w", assigneeUsername, teamID, err)
	}
	if !exists {
		return fmt.Errorf("user with username %q is not a member of team %d", assigneeUsername, teamID)
	}
	creatorMember, exists, err := members.MemberWithUsername(ctx, creatorUser.Username())
	if err != nil {
		return fmt.Errorf("getting member with username %q in team %d: %w", creatorUser.Username(), teamID, err)
	}
	if !exists {
		return fmt.Errorf("user with username %q is not a member of team %d", creatorUser.Username(), teamID)
	}
	_, err = creatorMember.CreateAssigned(ctx, title, description, points, assigneeMember.ID())
	if err != nil {
		return fmt.Errorf("creating assigned task in team %d: %w", teamID, err)
	}
	return nil
}

func CreateAssigned(community domain.Community, forms forms.Forms) core.Action {
	return createAssignedAction{
		community: community,
		forms:     forms,
	}
}

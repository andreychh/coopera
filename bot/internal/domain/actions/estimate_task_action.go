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

type estimateTaskAction struct {
	community domain.Community
	forms     forms.Forms
}

func (e estimateTaskAction) Perform(ctx context.Context, update telegram.Update) error {
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return fmt.Errorf("chat ID not found in update")
	}
	form := e.forms.Form(chatID)
	// Task ID
	taskIDStr, err := form.Field("task_id").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting task_id: %w", err)
	}
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("parsing task_id: %w", err)
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
	// Estimate Task
	user, exists, err := e.community.UserWithTelegramID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("getting user with telegram ID %d: %w", chatID, err)
	}
	if !exists {
		return fmt.Errorf("user with telegram ID %d does not exist", chatID)
	}
	task, exists, err := e.community.Task(ctx, taskID)
	if err != nil {
		return fmt.Errorf("getting task %d: %w", taskID, err)
	}
	if !exists {
		return fmt.Errorf("task %d does not exist", taskID)
	}
	team, err := task.Team(ctx)
	if err != nil {
		return fmt.Errorf("getting team for task %d: %w", task.ID(), err)
	}
	members, err := team.Members(ctx)
	if err != nil {
		return fmt.Errorf("getting members of team %d: %w", team.ID(), err)
	}
	member, exists, err := members.MemberWithUsername(ctx, user.Username())
	if err != nil {
		return fmt.Errorf("getting member with username %q in team %d: %w", user.Username(), team.ID(), err)
	}
	if !exists {
		return fmt.Errorf("member with username %q does not exist in team %d", user.Username(), team.ID())
	}
	err = member.EstimateTask(ctx, task.ID(), points)
	if err != nil {
		return fmt.Errorf("estimating task %d by member %q: %w", task.ID(), member.Username(), err)
	}
	return nil
}

func EstimateTask(community domain.Community, forms forms.Forms) core.Action {
	return estimateTaskAction{
		community: community,
		forms:     forms,
	}
}

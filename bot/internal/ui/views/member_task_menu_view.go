package views

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type memberTaskMenuView struct {
	community domain.Community
}

func (t memberTaskMenuView) Value(ctx context.Context, update telegram.Update) (content.Content, error) {
	callbackData, exists := attributes.CallbackData().Value(update)
	if !exists {
		return nil, fmt.Errorf("getting callback data from update: callback data not found")
	}
	taskID, err := protocol.ParseTaskID(callbackData)
	if err != nil {
		return nil, fmt.Errorf("parsing task ID from callback data %q: %w", callbackData, err)
	}
	task, exists, err := t.community.Task(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("getting task %d: %w", taskID, err)
	}
	if !exists {
		return nil, fmt.Errorf("task %d does not exist", taskID)
	}
	description, err := t.description(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("getting description for task %d: %w", taskID, err)
	}
	team, err := task.Team(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting team for task %d: %w", task.ID(), err)
	}
	if task.Status() == domain.StatusInProgress {
		return keyboards.Inline(
			content.Text(description),
			buttons.Matrix(
				buttons.Row(buttons.CallbackButton("Submit for review", protocol.ToMemberTaskMenu(taskID))),
				buttons.Row(buttons.CallbackButton("My tasks", protocol.ToMemberTasksMenu(team.ID()))),
			),
		), nil
	}
	return keyboards.Inline(
		content.Text(description),
		buttons.Matrix(
			buttons.Row(buttons.CallbackButton("My tasks", protocol.ToMemberTasksMenu(team.ID()))),
		),
	), nil
}

func (t memberTaskMenuView) description(ctx context.Context, task domain.Task) (string, error) {
	team, err := task.Team(ctx)
	if err != nil {
		return "", fmt.Errorf("getting team for task %d: %w", task.ID(), err)
	}
	points, exists := task.Points()
	if !exists {
		return "", fmt.Errorf("getting points for task %d: points not set", task.ID())
	}
	return fmt.Sprintf(
		"Task %q\nCreated in team %q\nAt %s\nDescription:\n%s\n\nPoints: %d | Status: %s\n",
		task.Title(),
		team.Name(),
		task.CreatedAt().Format("02.01.2006 15:04"),
		task.Description(),
		points,
		task.Status(),
	), nil
}

func MemberTaskMenuView(community domain.Community) sources.Source[content.Content] {
	return memberTaskMenuView{community: community}
}

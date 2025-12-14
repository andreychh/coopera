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

type userTaskMenuView struct {
	community domain.Community
}

func (t userTaskMenuView) Value(ctx context.Context, update telegram.Update) (content.Content, error) {
	callbackData, exists := attributes.CallbackData().Value(update)
	if !exists {
		return nil, fmt.Errorf("getting callback data from update: callback data not found")
	}
	id, err := protocol.ParseTaskID(callbackData)
	if err != nil {
		return nil, fmt.Errorf("parsing task ID from callback data %q: %w", callbackData, err)
	}
	task, err := t.community.Task(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting task %d: %w", id, err)
	}
	description, err := t.description(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("getting description for task %d: %w", id, err)
	}
	if task.Status() == domain.StatusAssigned {
		return keyboards.Inline(
			content.Text(description),
			buttons.Matrix(
				buttons.Row(buttons.CallbackButton("Mark as completed", protocol.ToUserTaskMenu(id))),
				buttons.Row(buttons.CallbackButton("Main menu", protocol.ToMainMenu())),
			),
		), nil
	}
	return keyboards.Inline(
		content.Text(description),
		buttons.Matrix(buttons.Row(buttons.CallbackButton("Main menu", protocol.ToMainMenu()))),
	), nil
}

func (t userTaskMenuView) description(ctx context.Context, task domain.Task) (string, error) {
	team, err := task.Team(ctx)
	if err != nil {
		return "", fmt.Errorf("getting team for task %d: %w", task.ID(), err)
	}
	creator, err := task.Creator(ctx)
	if err != nil {
		return "", fmt.Errorf("getting creator for task %d: %w", task.ID(), err)
	}
	createdAt, err := task.CreatedAt(ctx)
	if err != nil {
		return "", fmt.Errorf("getting creation time for task %d: %w", task.ID(), err)
	}
	return fmt.Sprintf(
		"Task %q\nCreated in team %q\nBy %s (%s)\nAt %s\nDescription:\n%s\n\nPoints: %d | Status: %s (Assigned to you)\n",
		task.Title(),
		team.Name(),
		creator.Username(),
		creator.Role(),
		createdAt.Format("02.01.2006 15:04"),
		task.Description(),
		task.Points(),
		task.Status(),
	), nil
}

func UserTaskMenuView(community domain.Community) sources.Source[content.Content] {
	return userTaskMenuView{community: community}
}

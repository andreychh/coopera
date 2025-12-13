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

type allMyTasksMenuView struct {
	community domain.Community
}

func (m allMyTasksMenuView) Value(ctx context.Context, update telegram.Update) (content.Content, error) {
	id, found := attributes.ChatID().Value(update)
	if !found {
		return nil, fmt.Errorf("chat ID not found in update")
	}
	user, err := m.community.UserWithTelegramID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting user with telegram ID %d: %w", id, err)
	}
	tasks, err := user.AssignedTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting assigned tasks for user %d: %w", id, err)
	}
	slice, err := tasks.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting assigned tasks slice for user %d: %w", id, err)
	}
	matrix, err := m.tasksMatrix(ctx, slice)
	if err != nil {
		return nil, fmt.Errorf("creating tasks matrix for user %d: %w", id, err)
	}
	return keyboards.Inline(
		content.Text("assigned tasks:"),
		matrix.WithRow(buttons.Row(buttons.CallbackButton("Main menu", protocol.ToMainMenu()))),
	), nil
}

func (m allMyTasksMenuView) tasksMatrix(ctx context.Context, tasks []domain.Task) (buttons.ButtonMatrix[buttons.InlineButton], error) {
	matrix := buttons.Matrix[buttons.InlineButton]()
	for _, task := range tasks {
		button, err := m.taskButton(ctx, task)
		if err != nil {
			return nil, fmt.Errorf("creating task button for task %d: %w", task.ID(), err)
		}
		matrix = matrix.WithRow(buttons.Row(button))
	}
	return matrix, nil
}

func (m allMyTasksMenuView) taskButton(ctx context.Context, task domain.Task) (buttons.InlineButton, error) {
	team, err := task.Team(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting team for task %d: %w", task.ID(), err)
	}
	return buttons.CallbackButton(
		fmt.Sprintf("%s | %s | %d | %s", team.Name(), task.Title(), task.Points(), task.Status()),
		"not_implemented",
		//protocol.ToTaskMenu(task.ID()),
	), nil
}

func TasksAssignedToUserView(community domain.Community) sources.Source[content.Content] {
	return allMyTasksMenuView{community: community}
}

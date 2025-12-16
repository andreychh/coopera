package views

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/formatting"
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
	user, exists, err := m.community.UserWithTelegramID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting user with telegram ID %d: %w", id, err)
	}
	if !exists {
		return nil, fmt.Errorf("user with telegram ID %d does not exist", id)
	}
	tasks, err := user.AssignedTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting assigned tasks for user %d: %w", id, err)
	}
	slice, err := tasks.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting assigned tasks slice for user %d: %w", id, err)
	}
	if len(slice) == 0 {
		text := `📋 <b>Мои задачи</b>
        
У вас пока нет активных задач ни в одной из команд.`
		return keyboards.Inline(
			formatting.Formatted(content.Text(text), formatting.ParseModeHTML),
			buttons.Matrix(
				buttons.Row(buttons.CallbackButton("🔙 Главное меню", protocol.ToMainMenu())),
			),
		), nil
	}
	matrix, err := m.tasksMatrix(ctx, slice)
	if err != nil {
		return nil, fmt.Errorf("creating tasks matrix for user %d: %w", id, err)
	}
	text := `📋 <b>Мои задачи</b>

Список задач, назначенных на вас во всех командах.

<b>Статусы:</b>
🔨 — В работе
👀 — На проверке
✅ — Выполнено`
	return keyboards.Inline(
		formatting.Formatted(content.Text(text), formatting.ParseModeHTML),
		matrix.WithRow(buttons.Row(buttons.CallbackButton("🔙 Главное меню", protocol.ToMainMenu()))),
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
	points, exists := task.Points()
	if !exists {
		points = 0
	}
	statusIcon := ""
	switch task.Status() {
	case domain.StatusInProgress:
		statusIcon = "🔨"
	case domain.StatusInReview:
		statusIcon = "👀"
	case domain.StatusDone:
		statusIcon = "✅"
	default:
		statusIcon = "❓"
	}
	label := fmt.Sprintf("%s %s: %s (+%d)", statusIcon, team.Name(), task.Title(), points)
	return buttons.CallbackButton(
		label,
		protocol.ToUserTaskMenu(task.ID()),
	), nil
}

func TasksAssignedToUserView(community domain.Community) sources.Source[content.Content] {
	return allMyTasksMenuView{community: community}
}

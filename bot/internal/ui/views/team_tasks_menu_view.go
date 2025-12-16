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

type allTeamTasksMenuView struct {
	community domain.Community
}

func (m allTeamTasksMenuView) Value(ctx context.Context, update telegram.Update) (content.Content, error) {
	callbackData, exists := attributes.CallbackData().Value(update)
	if !exists {
		return nil, fmt.Errorf("getting callback data from update: callback data not found")
	}
	id, err := protocol.ParseTeamID(callbackData)
	if err != nil {
		return nil, fmt.Errorf("parsing team ID from callback data %q: %w", callbackData, err)
	}
	team, exists, err := m.community.Team(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", id, err)
	}
	if !exists {
		return nil, fmt.Errorf("team %d does not exist", id)
	}
	tasks, err := team.Tasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting tasks for team %d: %w", id, err)
	}
	slice, err := tasks.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting tasks slice for team %d: %w", id, err)
	}
	if len(slice) == 0 {
		text := fmt.Sprintf(`📋 <b>Доска задач: %s</b>

В этой команде пока нет задач.
Вернитесь в меню команды, чтобы создать первую задачу.`, team.Name())
		return keyboards.Inline(
			formatting.Formatted(content.Text(text), formatting.ParseModeHTML),
			buttons.Matrix(
				buttons.Row(buttons.CallbackButton("🔙 Меню команды", protocol.ToTeamMenu(id))),
			),
		), nil
	}
	matrix, err := m.tasksMatrix(ctx, slice)
	if err != nil {
		return nil, fmt.Errorf("creating tasks matrix for team %d: %w", id, err)
	}
	text := fmt.Sprintf(`📋 <b>Доска задач: %s</b>

Используйте этот список для отслеживания хода выполнения задач команды.

<b>Статусы:</b>
📝 — Требует оценки
🗄 — Ожидает исполнителя
🔨 — В работе (с исполнителем)
👀 — На проверке
✅ — Завершено`, team.Name())
	return keyboards.Inline(
		formatting.Formatted(content.Text(text), formatting.ParseModeHTML),
		matrix.WithRow(buttons.Row(buttons.CallbackButton("🔙 Меню команды", protocol.ToTeamMenu(id)))),
	), nil
}

func (m allTeamTasksMenuView) tasksMatrix(ctx context.Context, tasks []domain.Task) (buttons.ButtonMatrix[buttons.InlineButton], error) {
	matrix := buttons.Matrix[buttons.InlineButton]()
	for _, task := range tasks {
		button, err := m.taskButton(ctx, task)
		if err != nil {
			return nil, fmt.Errorf("creating button for task %d: %w", task.ID(), err)
		}
		matrix = matrix.WithRow(buttons.Row(button))
	}
	return matrix, nil
}

func (m allTeamTasksMenuView) taskButton(ctx context.Context, task domain.Task) (buttons.InlineButton, error) {
	if task.Status() == domain.StatusDraft {
		return buttons.CallbackButton(
			fmt.Sprintf("📝 %s (Оценка...)", task.Title()),
			protocol.ToTeamTaskMenu(task.ID()),
		), nil
	}
	points, exists := task.Points()
	if !exists {
		points = 0
	}
	statusIcon := ""
	needsAssignee := false
	switch task.Status() {
	case domain.StatusOpen:
		statusIcon = "🗄"
	case domain.StatusInProgress:
		statusIcon = "🔨"
		needsAssignee = true
	case domain.StatusInReview:
		statusIcon = "👀"
		needsAssignee = true
	case domain.StatusDone:
		statusIcon = "✅"
		needsAssignee = true
	default:
		statusIcon = "❓"
	}
	assigneeStr := ""
	if needsAssignee {
		assignee, found, err := task.Assignee(ctx)
		if err != nil {
			return nil, fmt.Errorf("getting assignee for task %d: %w", task.ID(), err)
		}
		if found {
			assigneeStr = fmt.Sprintf(" @%s", assignee.Username())
		} else {
			assigneeStr = " (no user)"
		}
	}
	label := fmt.Sprintf("%s %s (+%d)%s", statusIcon, task.Title(), points, assigneeStr)
	return buttons.CallbackButton(
		label,
		protocol.ToTeamTaskMenu(task.ID()),
	), nil
}

func TeamTasks(community domain.Community) sources.Source[content.Content] {
	return allTeamTasksMenuView{community: community}
}

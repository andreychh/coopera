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
	team, err := task.Team(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting team for task %d: %w", task.ID(), err)
	}
	description, err := t.description(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("getting description for task %d: %w", taskID, err)
	}
	btns := buttons.Matrix[buttons.InlineButton]()
	if task.Status() == domain.StatusInProgress {
		btns = btns.WithRow(buttons.Row(
			buttons.CallbackButton("📤 Отправить на проверку", protocol.ToMemberTaskMenu(taskID)),
		))
	}
	btns = btns.WithRow(buttons.Row(
		buttons.CallbackButton("🔙 К задачам команды", protocol.ToMemberTasksMenu(team.ID())),
	))
	return keyboards.Inline(
		formatting.Formatted(content.Text(description), formatting.ParseModeHTML),
		btns,
	), nil
}

func (t memberTaskMenuView) description(ctx context.Context, task domain.Task) (string, error) {
	team, err := task.Team(ctx)
	if err != nil {
		return "", fmt.Errorf("getting team for task %d: %w", task.ID(), err)
	}
	points, exists := task.Points()
	if !exists {
		points = 0
	}
	creator, err := task.CreatedBy(ctx)
	username := "unknown"
	if err == nil {
		username = creator.Username()
	}
	statusStr := ""
	switch task.Status() {
	case domain.StatusInProgress:
		statusStr = "🔨 В работе"
	case domain.StatusInReview:
		statusStr = "👀 На проверке"
	case domain.StatusDone:
		statusStr = "✅ Выполнено"
	default:
		statusStr = string(task.Status())
	}
	return fmt.Sprintf(`📄 <b>Задача: %s</b>

<b>Команда:</b> %s
<b>Автор:</b> @%s
<b>Создана:</b> %s

<b>Статус:</b> %s
<b>Стоимость:</b> +%d баллов

<b>Описание:</b>
<i>%s</i>`,
		task.Title(),
		team.Name(),
		username,
		task.CreatedAt().Format("02.01.2006 15:04"),
		statusStr,
		points,
		task.Description(),
	), nil
}

func MemberTaskMenuView(community domain.Community) sources.Source[content.Content] {
	return memberTaskMenuView{community: community}
}

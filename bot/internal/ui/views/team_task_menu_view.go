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

type teamTaskMenuView struct {
	community domain.Community
}

func (t teamTaskMenuView) Value(ctx context.Context, update telegram.Update) (content.Content, error) {
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
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return nil, fmt.Errorf("chat ID not found in update")
	}
	currentUser, exists, err := t.community.UserWithTelegramID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("getting user with telegram ID %d: %w", chatID, err)
	}
	if !exists {
		return nil, fmt.Errorf("user with telegram ID %d does not exist", chatID)
	}
	members, err := team.Members(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting members of team %d: %w", team.ID(), err)
	}
	currentMember, exists, err := members.MemberWithUsername(ctx, currentUser.Username())
	if err != nil {
		return nil, fmt.Errorf("getting member for user %d in team %d: %w", currentUser.ID(), team.ID(), err)
	}
	if !exists {
		return nil, fmt.Errorf("member for user %d in team %d does not exist", currentUser.ID(), team.ID())
	}
	assigneeMember, assigneeFound, err := task.Assignee(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting assignee member for task %d: %w", task.ID(), err)
	}
	isAssignee := assigneeFound && assigneeMember.ID() == currentMember.ID()
	description, err := t.description(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("getting description for task %d: %w", taskID, err)
	}
	btns := buttons.Matrix[buttons.InlineButton]()
	if task.Status() == domain.StatusDraft && currentMember.Role() == domain.RoleManager {
		btns = btns.WithRow(buttons.Row(
			buttons.CallbackButton("⭐️ Оценить задачу", protocol.StartEstimateTaskForm(task.ID())),
		))
	} else if task.Status() == domain.StatusOpen {
		btns = btns.WithRow(buttons.Row(
			buttons.CallbackButton("🙋‍♂️ Взять в работу", protocol.ToTeamTaskMenuWithAction(task.ID(), protocol.ActionAssignTaskToSelf)),
		))
	} else if task.Status() == domain.StatusInProgress && isAssignee {
		btns = btns.WithRow(buttons.Row(
			buttons.CallbackButton("📤 Отправить на проверку", protocol.ToTeamTaskMenuWithAction(task.ID(), protocol.ActionSubmitTaskForReview)),
		))
	} else if task.Status() == domain.StatusInReview && currentMember.Role() == domain.RoleManager {
		btns = btns.WithRow(buttons.Row(
			buttons.CallbackButton("✅ Подтвердить выполнение", protocol.ToTeamTaskMenuWithAction(task.ID(), protocol.ActionApproveTask)),
		))
	}
	btns = btns.WithRow(buttons.Row(
		buttons.CallbackButton("🔙 К доске задач", protocol.ToTeamTasksMenu(team.ID())),
	))
	return keyboards.Inline(
		formatting.Formatted(content.Text(description), formatting.ParseModeHTML),
		btns,
	), nil
}

func (t teamTaskMenuView) description(ctx context.Context, task domain.Task) (string, error) {
	team, err := task.Team(ctx)
	if err != nil {
		return "", fmt.Errorf("getting team for task %d: %w", task.ID(), err)
	}
	creator, err := task.CreatedBy(ctx)
	username := "unknown"
	if err == nil {
		username = creator.Username()
	}
	assigneeStr := ""
	assignee, found, err := task.Assignee(ctx)
	if err == nil && found {
		assigneeStr = fmt.Sprintf("\n<b>Исполнитель:</b> @%s", assignee.Username())
	}
	pointsStr := "<i>(требует оценки)</i>"
	if p, exists := task.Points(); exists {
		pointsStr = fmt.Sprintf("+%d баллов", p)
	}
	statusStr := ""
	switch task.Status() {
	case domain.StatusDraft:
		statusStr = "📝 Требует оценки"
	case domain.StatusOpen:
		statusStr = "🗄 Открыта (Backlog)"
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
<b>Автор:</b> @%s%s
<b>Создана:</b> %s

<b>Статус:</b> %s
<b>Стоимость:</b> %s

<b>Описание:</b>
<i>%s</i>`,
		task.Title(),
		team.Name(),
		username,
		assigneeStr,
		task.CreatedAt().Format("02.01.2006 15:04"),
		statusStr,
		pointsStr,
		task.Description(),
	), nil
}

func TeamTaskMenuView(community domain.Community) sources.Source[content.Content] {
	return teamTaskMenuView{community: community}
}

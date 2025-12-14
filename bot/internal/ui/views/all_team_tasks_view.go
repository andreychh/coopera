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
	matrix, err := m.tasksMatrix(ctx, slice)
	if err != nil {
		return nil, fmt.Errorf("creating tasks matrix for team %d: %w", id, err)
	}
	return keyboards.Inline(
		content.Text(fmt.Sprintf("Team %s tasks:", team.Name())),
		matrix.WithRow(buttons.Row(buttons.CallbackButton("Team menu", protocol.ToTeamMenu(id)))),
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
			fmt.Sprintf("%q | Open (needs estimation)", task.Title()),
			"not_implemented",
		), nil
	}
	points, exists := task.Points()
	if !exists {
		return nil, fmt.Errorf("task %d has status %q but no points assigned", task.ID(), task.Status())
	}
	var statusLabel string
	needsAssignee := false
	switch task.Status() {
	case domain.StatusOpen:
		statusLabel = "Open"
	case domain.StatusInProgress:
		statusLabel = "In Progress"
		needsAssignee = true
	case domain.StatusInReview:
		statusLabel = "In Review"
		needsAssignee = true
	case domain.StatusDone:
		statusLabel = "Done"
		needsAssignee = true
	default:
		return nil, fmt.Errorf("unknown task status %q for task %d", task.Status(), task.ID())
	}
	var assigneeStr string
	if needsAssignee {
		assignee, found, err := task.Assignee(ctx)
		if err != nil {
			return nil, fmt.Errorf("getting assignee for task %d: %w", task.ID(), err)
		}
		if !found {
			return nil, fmt.Errorf("task %d is %q but has no assignee", task.ID(), statusLabel)
		}
		assigneeStr = fmt.Sprintf(" (@%s)", assignee.Username())
	}
	return buttons.CallbackButton(
		fmt.Sprintf("%q | %d | %s%s", task.Title(), points, statusLabel, assigneeStr),
		"not_implemented", // protocol.ToTaskMenu(task.ID()),
	), nil
}

func AllTeamTasks(community domain.Community) sources.Source[content.Content] {
	return allTeamTasksMenuView{community: community}
}

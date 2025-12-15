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
	description, err := t.description(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("getting description for task %d: %w", taskID, err)
	}
	team, err := task.Team(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting team for task %d: %w", task.ID(), err)
	}
	members, err := team.Members(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting members of team %d: %w", team.ID(), err)
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
	currentMember, exists, err := members.MemberWithUsername(ctx, currentUser.Username())
	if err != nil {
		return nil, fmt.Errorf("getting member for user %d in team %d: %w", currentUser.ID(), team.ID(), err)
	}
	if !exists {
		return nil, fmt.Errorf("member for user %d in team %d does not exist", currentUser.ID(), team.ID())
	}
	assigneeMember, _, err := task.Assignee(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting assignee member for task %d: %w", task.ID(), err)
	}
	if task.Status() == domain.StatusDraft && currentMember.Role() == domain.RoleManager {
		return keyboards.Inline(
			content.Text(description),
			buttons.Matrix(
				buttons.Row(buttons.CallbackButton(
					"Estimate",
					protocol.StartEstimateTaskForm(task.ID()),
				)),
				buttons.Row(buttons.CallbackButton(
					"Team tasks",
					protocol.ToTeamTasksMenu(team.ID()),
				)),
			),
		), nil
	} else if task.Status() == domain.StatusOpen {
		return keyboards.Inline(
			content.Text(description),
			buttons.Matrix(
				buttons.Row(buttons.CallbackButton(
					"Take task",
					protocol.ToTeamTaskMenuWithAction(task.ID(), protocol.ActionAssignTaskToSelf),
				)),
				buttons.Row(buttons.CallbackButton(
					"Team tasks",
					protocol.ToTeamTasksMenu(team.ID()),
				)),
			),
		), nil
	} else if task.Status() == domain.StatusInProgress && currentMember.ID() == assigneeMember.ID() {
		return keyboards.Inline(
			content.Text(description),
			buttons.Matrix(
				buttons.Row(buttons.CallbackButton(
					"Submit task for review",
					protocol.ToTeamTaskMenuWithAction(task.ID(), protocol.ActionSubmitTaskForReview),
				)),
				buttons.Row(buttons.CallbackButton(
					"Team tasks",
					protocol.ToTeamTasksMenu(team.ID()),
				)),
			),
		), nil
	} else if task.Status() == domain.StatusInReview && currentMember.Role() == domain.RoleManager {
		return keyboards.Inline(
			content.Text(description),
			buttons.Matrix(
				buttons.Row(buttons.CallbackButton(
					"Approve task",
					protocol.ToTeamTaskMenuWithAction(task.ID(), protocol.ActionApproveTask),
				)),
				buttons.Row(buttons.CallbackButton(
					"Team tasks",
					protocol.ToTeamTasksMenu(team.ID()),
				)),
			),
		), nil
	}
	return keyboards.Inline(
		content.Text(description),
		buttons.Matrix(
			buttons.Row(buttons.CallbackButton("Team tasks", protocol.ToTeamTasksMenu(team.ID()))),
		),
	), nil
}

func (t teamTaskMenuView) description(ctx context.Context, task domain.Task) (string, error) {
	team, err := task.Team(ctx)
	if err != nil {
		return "", fmt.Errorf("getting team for task %d: %w", task.ID(), err)
	}
	points, exists := task.Points()
	if !exists {
		return fmt.Sprintf(
			"Task %q\nCreated in team %q\nAt %s\nDescription:\n%s\n\nPoints: (needs estimation) | Status: %s\n",
			task.Title(),
			team.Name(),
			task.CreatedAt().Format("02.01.2006 15:04"),
			task.Description(),
			task.Status(),
		), nil
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

func TeamTaskMenuView(community domain.Community) sources.Source[content.Content] {
	return teamTaskMenuView{community: community}
}

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

type memberTasksMenuView struct {
	community domain.Community
}

func (m memberTasksMenuView) Value(ctx context.Context, update telegram.Update) (content.Content, error) {
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return nil, fmt.Errorf("chat ID not found in update")
	}
	callbackData, exists := attributes.CallbackData().Value(update)
	if !exists {
		return nil, fmt.Errorf("getting callback data from update: callback data not found")
	}
	teamID, err := protocol.ParseTeamID(callbackData)
	if err != nil {
		return nil, fmt.Errorf("parsing team ID from callback data %q: %w", callbackData, err)
	}
	team, exists, err := m.community.Team(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", teamID, err)
	}
	if !exists {
		return nil, fmt.Errorf("team %d does not exist", teamID)
	}
	user, exists, err := m.community.UserWithTelegramID(ctx, chatID)
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
	member, exists, err := members.MemberWithUsername(ctx, user.Username())
	if err != nil {
		return nil, fmt.Errorf("getting member for user %d in team %d: %w", user.ID(), team.ID(), err)
	}
	if !exists {
		return nil, fmt.Errorf("member for user %d in team %d does not exist", user.ID(), team.ID())
	}
	tasks, err := member.AssignedTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting assigned tasks for user %d: %w", user.ID(), err)
	}
	slice, err := tasks.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting assigned tasks slice for user %d: %w", user.ID(), err)
	}
	matrix, err := m.tasksMatrix(ctx, slice)
	if err != nil {
		return nil, fmt.Errorf("creating tasks matrix for user %d: %w", user.ID(), err)
	}
	return keyboards.Inline(
		content.Text("your tasks:"),
		matrix.
			WithRow(buttons.Row(buttons.CallbackButton("Team menu", protocol.ToTeamMenu(team.ID())))),
	), nil
}

func (m memberTasksMenuView) tasksMatrix(ctx context.Context, tasks []domain.Task) (buttons.ButtonMatrix[buttons.InlineButton], error) {
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

func (m memberTasksMenuView) taskButton(_ context.Context, task domain.Task) (buttons.InlineButton, error) {
	points, exists := task.Points()
	if !exists {
		return nil, fmt.Errorf("getting points for task %d: points not set", task.ID())
	}
	return buttons.CallbackButton(
		fmt.Sprintf("%q | %d | %s", task.Title(), points, task.Status()),
		protocol.ToMemberTaskMenu(task.ID()),
	), nil
}

func MemberTasksMenuView(community domain.Community) sources.Source[content.Content] {
	return memberTasksMenuView{community: community}
}

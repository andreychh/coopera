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
	team, err := m.community.Team(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", teamID, err)
	}
	user, err := m.community.UserWithTelegramID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("getting user with telegram ID %d: %w", chatID, err)
	}
	member, err := team.MemberWithUserID(ctx, user.ID())
	if err != nil {
		return nil, fmt.Errorf("getting member with user ID %d in team %d: %w", user.ID(), teamID, err)
	}
	tasks, err := member.Tasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting assigned tasks for user %d: %w", user.ID(), err)
	}
	slice, err := tasks.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting assigned tasks slice for user %d: %w", user.ID(), err)
	}
	return keyboards.Inline(
		content.Text("your tasks:"),
		m.tasksMatrix(slice).
			WithRow(buttons.Row(buttons.CallbackButton("Team menu", protocol.ToTeamMenu(team.ID())))),
	), nil
}

func (m memberTasksMenuView) tasksMatrix(tasks []domain.Task) buttons.ButtonMatrix[buttons.InlineButton] {
	matrix := buttons.Matrix[buttons.InlineButton]()
	for _, task := range tasks {
		matrix = matrix.WithRow(buttons.Row(m.taskButton(task)))
	}
	return matrix
}

func (m memberTasksMenuView) taskButton(task domain.Task) buttons.InlineButton {
	return buttons.CallbackButton(
		fmt.Sprintf("%s | %d | %s", task.Title(), task.Points(), task.Status()),
		"not_implemented",
		// protocol.ToTaskMenu(task.ID()),
	)
}

func MemberTasksMenuView(community domain.Community) sources.Source[content.Content] {
	return memberTasksMenuView{community: community}
}

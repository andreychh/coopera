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
	team, err := m.community.Team(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", id, err)
	}
	tasks, err := team.Tasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting tasks for team %d: %w", id, err)
	}
	slice, err := tasks.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting tasks slice for team %d: %w", id, err)
	}
	return keyboards.Inline(
		content.Text(fmt.Sprintf("Team %s tasks:", team.Name())),
		m.tasksMatrix(slice).WithRow(buttons.Row(buttons.CallbackButton("Team menu", protocol.ToTeamMenu(id)))),
	), nil
}

func (m allTeamTasksMenuView) tasksMatrix(tasks []domain.Task) buttons.ButtonMatrix[buttons.InlineButton] {
	matrix := buttons.Matrix[buttons.InlineButton]()
	for _, task := range tasks {
		matrix = matrix.WithRow(buttons.Row(m.taskButton(task)))
	}
	return matrix
}

func (m allTeamTasksMenuView) taskButton(task domain.Task) buttons.InlineButton {
	return buttons.CallbackButton(
		fmt.Sprintf("%s | %d | %s", task.Title(), task.Points(), task.Status()),
		"not_implemented",
		// protocol.ToTaskMenu(task.ID()),
	)
}

func AllTeamTasks(community domain.Community) sources.Source[content.Content] {
	return allTeamTasksMenuView{community: community}
}

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

type teamMenuView struct {
	community domain.Community
}

func (t teamMenuView) Value(ctx context.Context, update telegram.Update) (content.Content, error) {
	callbackData, exists := attributes.CallbackData().Value(update)
	if !exists {
		return nil, fmt.Errorf("getting callback data from update: callback data not found")
	}
	id, err := protocol.ParseTeamID(callbackData)
	if err != nil {
		return nil, fmt.Errorf("parsing team ID from callback data %q: %w", callbackData, err)
	}
	team, exists, err := t.community.Team(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", id, err)
	}
	if !exists {
		return nil, fmt.Errorf("team %d does not exist", id)
	}
	stats, err := team.Stats(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting stats for team %d: %w", team.ID(), err)
	}
	text := fmt.Sprintf(`🏢 <b>Команда: %s</b>

📊 <b>Сводка по задачам:</b>

<b>Очередь:</b> %d шт. (+%d)
<b>В работе:</b> %d шт. (+%d)
<b>На проверке:</b> %d шт. (+%d)
<b>Завершено:</b> %d шт. (+%d)

👇 Выберите действие:`,
		team.Name(),
		stats.Backlog.UnassignedCount, stats.Backlog.UnassignedPoints,
		stats.ActiveWork.InProgressCount, stats.ActiveWork.InProgressPoints,
		stats.ActiveWork.InReviewCount, stats.ActiveWork.InReviewPoints,
		stats.Achievements.CompletedCount, stats.Achievements.CompletedPoints,
	)

	return keyboards.Inline(
		formatting.Formatted(content.Text(text), formatting.ParseModeHTML),
		buttons.Matrix(
			buttons.Row(buttons.CallbackButton("➕ Создать задачу", protocol.StartCreateTaskForm(team.ID()))),
			buttons.Row(buttons.CallbackButton("👤 Мои задачи", protocol.ToMemberTasksMenu(team.ID()))),
			buttons.Row(buttons.CallbackButton("📋 Доска задач", protocol.ToTeamTasksMenu(team.ID()))),
			buttons.Row(buttons.CallbackButton("👥 Участники", protocol.ToMembersMenu(team.ID()))),
			buttons.Row(buttons.CallbackButton("🔙 К списку команд", protocol.ToTeamsMenu())),
		),
	), nil
}

func TeamMenu(community domain.Community) sources.Source[content.Content] {
	return teamMenuView{community: community}
}

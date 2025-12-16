package views

import (
	"context"
	"fmt"
	"sort"
	"strings"

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

type userStatsMenuView struct {
	community domain.Community
}

func (u userStatsMenuView) Value(ctx context.Context, update telegram.Update) (content.Content, error) {
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return nil, fmt.Errorf("chat ID not found in update")
	}
	user, exists, err := u.community.UserWithTelegramID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("getting user with telegram ID %d: %w", chatID, err)
	}
	if !exists {
		return nil, fmt.Errorf("user with telegram ID %d does not exist", chatID)
	}
	stats, err := user.Stats(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting stats for user %d: %w", user.ID(), err)
	}
	if len(stats.Teams) == 0 {
		text := `📊 <b>Личная статистика</b>

Вы пока не состоите ни в одной команде, поэтому статистика пуста.
Вступите в команду или создайте свою, чтобы начать трекинг.`
		return keyboards.Inline(
			formatting.Formatted(content.Text(text), formatting.ParseModeHTML),
			buttons.Matrix(
				buttons.Row(buttons.CallbackButton("🔙 Главное меню", protocol.ToMainMenu())),
			),
		), nil
	}
	var sb strings.Builder
	sb.WriteString("📊 <b>Личная статистика</b>\n\n")
	sb.WriteString("Сводка вашей эффективности по всем командам.\n\n")
	teamNames := make([]string, 0, len(stats.Teams))
	for name := range stats.Teams {
		teamNames = append(teamNames, name)
	}
	sort.Strings(teamNames)
	for _, name := range teamNames {
		tStat := stats.Teams[name]
		sb.WriteString(fmt.Sprintf("🏢 <b>%s</b>\n", name))
		if tStat.ActiveLoad.TasksCount > 0 {
			sb.WriteString(fmt.Sprintf("<b>В работе:</b> %d шт. (+%d)\n",
				tStat.ActiveLoad.TasksCount, tStat.ActiveLoad.TotalPoints))
		} else {
			sb.WriteString("<b>В работе:</b> нет задач\n")
		}
		sb.WriteString(fmt.Sprintf("<b>Завершено:</b> %d шт. (+%d)\n\n",
			tStat.LifetimeContribution.TasksCompleted, tStat.LifetimeContribution.PointsEarned))
	}
	return keyboards.Inline(
		formatting.Formatted(content.Text(sb.String()), formatting.ParseModeHTML),
		buttons.Matrix(
			buttons.Row(buttons.CallbackButton("🔙 Главное меню", protocol.ToMainMenu())),
		),
	), nil
}

func UserStatsMenu(community domain.Community) sources.Source[content.Content] {
	return userStatsMenuView{
		community: community,
	}
}

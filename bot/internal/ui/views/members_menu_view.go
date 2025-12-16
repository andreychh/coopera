package views

import (
	"context"
	"fmt"
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

type membersMenuView struct {
	community domain.Community
}

func (m membersMenuView) Value(ctx context.Context, update telegram.Update) (content.Content, error) {
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
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return nil, fmt.Errorf("chat ID not found")
	}
	user, exists, err := m.community.UserWithTelegramID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("getting user %d: %w", chatID, err)
	}
	if !exists {
		return nil, fmt.Errorf("user %d not found", chatID)
	}
	members, err := team.Members(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting members: %w", err)
	}
	currentMember, exists, err := members.MemberWithUsername(ctx, user.Username())
	if err != nil {
		return nil, fmt.Errorf("getting current member: %w", err)
	}
	isManager := exists && currentMember.Role() == domain.RoleManager
	text, err := m.renderMembersList(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("generating members text: %w", err)
	}
	btns := buttons.Matrix[buttons.InlineButton]()
	if isManager {
		btns = btns.WithRow(buttons.Row(
			buttons.CallbackButton("➕ Добавить участника", protocol.StartAddMemberForm(team.ID())),
		))
	}
	btns = btns.WithRow(buttons.Row(
		buttons.CallbackButton("🔙 Меню команды", protocol.ToTeamMenu(team.ID())),
	))
	return keyboards.Inline(
		formatting.Formatted(content.Text(text), formatting.ParseModeHTML),
		btns,
	), nil
}

func (m membersMenuView) renderMembersList(ctx context.Context, team domain.Team) (string, error) {
	members, err := team.Members(ctx)
	if err != nil {
		return "", err
	}
	slice, err := members.All(ctx)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("👥 <b>Участники команды: %s</b>\n\n", team.Name()))
	sb.WriteString("Список участников и их вклад в команде:\n\n")
	for _, member := range slice {
		stats, err := member.Stats(ctx)
		if err != nil {
			stats = domain.MemberStats{}
		}
		roleIcon := "👤"
		if member.Role() == domain.RoleManager {
			roleIcon = "⭐"
		}
		sb.WriteString(fmt.Sprintf("%s <b>@%s</b>\n", roleIcon, member.Username()))
		sb.WriteString(fmt.Sprintf("<b>В работе:</b> %d (+%d)\n",
			stats.CurrentState.AssignedTasks, stats.CurrentState.AssignedPoints,
		))
		sb.WriteString(fmt.Sprintf("<b>Завершено:</b> %d (+%d)\n\n",
			stats.Contribution.TasksDone, stats.Contribution.PointsDone,
		))
	}
	return sb.String(), nil
}

func MembersMenu(community domain.Community) sources.Source[content.Content] {
	return membersMenuView{community: community}
}

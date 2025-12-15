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
	members, err := team.Members(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting members for team %d: %w", teamID, err)
	}
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return nil, fmt.Errorf("chat ID not found in update")
	}
	user, exists, err := m.community.UserWithTelegramID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("getting user with telegram ID %d: %w", chatID, err)
	}
	if !exists {
		return nil, fmt.Errorf("user with telegram ID %d does not exist", chatID)
	}
	member, exists, err := members.MemberWithUsername(ctx, user.Username())
	if err != nil {
		return nil, fmt.Errorf("getting member with user ID %d in team %d: %w", user.ID(), teamID, err)
	}
	if !exists {
		return nil, fmt.Errorf("member with user ID %d in team %d does not exist", user.ID(), teamID)
	}
	text, err := m.membersText(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("generating members text: %w", err)
	}
	if member.Role() == domain.RoleManager {
		return keyboards.Inline(
			content.Text(text),
			buttons.Matrix(
				buttons.Row(buttons.CallbackButton("Add member", protocol.StartAddMemberForm(team.ID()))),
				buttons.Row(buttons.CallbackButton("Team menu", protocol.ToTeamMenu(team.ID()))),
			),
		), nil
	}
	return keyboards.Inline(
		content.Text(text),
		buttons.Matrix(
			buttons.Row(buttons.CallbackButton("Team menu", protocol.ToTeamMenu(team.ID()))),
		),
	), nil
}

func (m membersMenuView) membersText(ctx context.Context, team domain.Team) (string, error) {
	members, err := team.Members(ctx)
	if err != nil {
		return "", fmt.Errorf("getting members for team %d: %w", team.ID(), err)
	}
	slice, err := members.All(ctx)
	if err != nil {
		return "", fmt.Errorf("getting all members for team %d: %w", team.ID(), err)
	}
	text := fmt.Sprintf("👥 *Team %s Members:*\n\n", team.Name())
	for _, member := range slice {
		memberText, err := m.memberText(ctx, member)
		if err != nil {
			return "", fmt.Errorf("generating text for member %d: %w", member.ID(), err)
		}
		text += memberText + "\n"
	}
	return text, nil
}

func (m membersMenuView) memberText(ctx context.Context, member domain.Member) (string, error) {
	if member.Role() == domain.RoleManager {
		return fmt.Sprintf("* @%s", member.Username()), nil
	}
	return fmt.Sprintf("  @%s", member.Username()), nil
}

func MembersMenu(community domain.Community) sources.Source[content.Content] {
	return membersMenuView{community: community}
}

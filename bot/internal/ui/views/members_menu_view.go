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
	slice, err := members.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting memebers slice: %w", err)
	}
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return nil, fmt.Errorf("chat ID not found in update")
	}
	user, exists, err := m.community.UserWithTelegramID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("getting user with telegram ID %d: %w", chatID, err)
	}
	member, exists, err := members.MemberWithUsername(ctx, user.Username())
	if err != nil {
		return nil, fmt.Errorf("getting member with user ID %d in team %d: %w", user.ID(), teamID, err)
	}
	if !exists {
		return nil, fmt.Errorf("member with user ID %d in team %d does not exist", user.ID(), teamID)
	}
	if member.Role() == domain.RoleManager {
		return keyboards.Inline(
			content.Text(fmt.Sprintf("Team %s members (As manager):", team.Name())),
			m.membersMatrix(slice).
				WithRow(buttons.Row(buttons.CallbackButton("Add member", protocol.StartAddMemberForm(team.ID())))).
				WithRow(buttons.Row(buttons.CallbackButton("Team menu", protocol.ToTeamMenu(team.ID())))),
		), nil
	}
	return keyboards.Inline(
		content.Text(fmt.Sprintf("Team %s members (As member):", team.Name())),
		m.membersMatrix(slice).
			WithRow(buttons.Row(buttons.CallbackButton("Team menu", protocol.ToTeamMenu(team.ID())))),
	), nil
}

func (m membersMenuView) membersMatrix(members []domain.Member) buttons.ButtonMatrix[buttons.InlineButton] {
	matrix := buttons.Matrix[buttons.InlineButton]()
	for _, member := range members {
		matrix = matrix.WithRow(buttons.Row(m.memberButton(member)))
	}
	return matrix
}

func (m membersMenuView) memberButton(member domain.Member) buttons.InlineButton {
	var text string
	if member.Role() == domain.RoleManager {
		text = fmt.Sprintf("* @%s", member.Username())
	} else {
		text = fmt.Sprintf("  @%s", member.Username())
	}
	return buttons.CallbackButton(text, protocol.ToMemberMenu(member.ID()))
}

func MembersMenu(community domain.Community) sources.Source[content.Content] {
	return membersMenuView{community: community}
}

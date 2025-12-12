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
	id, err := protocol.ParseTeamID(callbackData)
	if err != nil {
		return nil, fmt.Errorf("parsing team ID from callback data %q: %w", callbackData, err)
	}
	team, err := m.community.Team(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", id, err)
	}
	members, err := team.Members(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting members for team %d: %w", id, err)
	}
	details, err := members.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting memebers slice: %w", err)
	}
	return keyboards.Inline(
		content.Text(fmt.Sprintf("Team %s members:", team.Name())),
		m.membersMatrix(details).
			WithRow(buttons.Row(buttons.CallbackButton("Add member", "not_implemented"))).
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
	return buttons.CallbackButton(member.Name(), protocol.ToMemberMenu(member.ID()))
}

func MembersMenu(community domain.Community) sources.Source[content.Content] {
	return membersMenuView{community: community}
}

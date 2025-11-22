package views

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type membersMenuView struct {
	community domain.Community
}

func (m membersMenuView) Render(ctx context.Context, update telegram.Update) (content.Content, error) {
	callbackData, exists := attrs.CallbackData(update).Value()
	if !exists {
		return nil, fmt.Errorf("getting callback data from update: callback data not found")
	}
	id, err := protocol.Navigation.ParseTeamID(callbackData)
	if err != nil {
		return nil, fmt.Errorf("parsing team ID from callback data %q: %w", callbackData, err)
	}
	members, err := m.community.Team(id).Members(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting members details for team %d: %w", id, err)
	}
	details, err := members_{members: members}.details(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting members details: %w", err)
	}
	return keyboards.Inline(
		content.Text("Team members:"),
		m.membersMatrix(details).WithRow(buttons.Row(buttons.CallbackButton(
			"Invite member",
			"not_implemented",
		))),
	), nil
}

func (m membersMenuView) membersMatrix(members []domain.MemberDetails) buttons.ButtonMatrix[buttons.InlineButton] {
	matrix := buttons.Matrix[buttons.InlineButton]()
	for _, member := range members {
		matrix = matrix.WithRow(buttons.Row(m.memberButton(member)))
	}
	return matrix
}

func (m membersMenuView) memberButton(member domain.MemberDetails) buttons.InlineButton {
	return buttons.CallbackButton(
		member.Name(),
		callbacks.OutcomingData("change_menu").
			With("menu_name", "member").
			With("member_id", strconv.FormatInt(member.ID(), 10)).
			String(),
	)
}

type members_ struct {
	members []domain.Member
}

func (m members_) details(ctx context.Context) ([]domain.MemberDetails, error) {
	var details []domain.MemberDetails
	for _, member := range m.members {
		detail, err := member.Details(ctx)
		if err != nil {
			return nil, fmt.Errorf("getting details for member: %w", err)
		}
		details = append(details, detail)
	}
	return details, nil
}

func MembersMenu(community domain.Community) content.View {
	return membersMenuView{community: community}
}

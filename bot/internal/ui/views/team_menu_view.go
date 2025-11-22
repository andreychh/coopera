package views

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type teamMenuView struct {
	community domain.Community
}

func (t teamMenuView) Render(ctx context.Context, update telegram.Update) (content.Content, error) {
	callbackData, exists := attrs.CallbackData(update).Value()
	if !exists {
		return nil, fmt.Errorf("getting callback data from update: callback data not found")
	}
	id, err := protocol.Navigation.ParseTeamID(callbackData)
	if err != nil {
		return nil, fmt.Errorf("parsing team ID from callback data %q: %w", callbackData, err)
	}
	details, err := t.community.Team(id).Details(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting details for team %d: %w", id, err)
	}
	return keyboards.Inline(
		content.Text(fmt.Sprintf("%s menu:", details.Name())),
		buttons.Matrix(
			buttons.Row(buttons.CallbackButton("Members", protocol.Navigation.ToTeamMembers(details.ID()))),
		),
	), nil
}

func TeamMenu(community domain.Community) content.View {
	return teamMenuView{community: community}
}

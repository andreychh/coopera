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
	team, err := t.community.Team(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting team %d: %w", id, err)
	}
	return keyboards.Inline(
		content.Text(fmt.Sprintf("Team %s:", team.Name())),
		buttons.Matrix(
			buttons.Row(buttons.CallbackButton("Members", protocol.ToMembersMenu(team.ID()))),
			buttons.Row(
				buttons.CallbackButton("All tasks", "not_implemented"),
				buttons.CallbackButton("My tasks", "not_implemented"),
			),
			buttons.Row(buttons.CallbackButton("Add task", "not_implemented")),
			buttons.Row(buttons.CallbackButton("Teams menu", protocol.ToTeamsMenu())),
		),
	), nil
}

func TeamMenu(community domain.Community) sources.Source[content.Content] {
	return teamMenuView{community: community}
}

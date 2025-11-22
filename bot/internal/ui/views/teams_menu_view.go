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

type teamsMenuView struct {
	community domain.Community
}

func (t teamsMenuView) Render(ctx context.Context, update telegram.Update) (content.Content, error) {
	chatID, exists := attrs.ChatID(update).Value()
	if !exists {
		return nil, fmt.Errorf("getting chat ID from update: chat ID not found")
	}
	teams, err := t.community.UserWithTelegramID(chatID).CreatedTeams(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting user's created teams: %w", err)
	}
	details, err := teams_{teams: teams}.details(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting teams details: %w", err)
	}
	return keyboards.Inline(
		content.Text("Select a team:"),
		t.teamsMatrix(details).
			WithRow(buttons.Row(buttons.CallbackButton("Create team", protocol.Forms.StartPayload("create_team")))).
			WithRow(buttons.Row(buttons.CallbackButton("Main menu", protocol.Navigation.Payload("main")))),
	), nil
}

func (t teamsMenuView) teamsMatrix(teams []domain.TeamDetails) buttons.ButtonMatrix[buttons.InlineButton] {
	matrix := buttons.Matrix[buttons.InlineButton]()
	for _, team := range teams {
		matrix = matrix.WithRow(buttons.Row(t.teamButton(team)))
	}
	return matrix
}

func (t teamsMenuView) teamButton(team domain.TeamDetails) buttons.InlineButton {
	return buttons.CallbackButton(team.Name(), protocol.Navigation.ToTeam(team.ID()))
}

// TODO: add Teams interface in domain package with Details method
type teams_ struct {
	teams []domain.Team
}

func (t teams_) details(ctx context.Context) ([]domain.TeamDetails, error) {
	var details []domain.TeamDetails
	for _, team := range t.teams {
		detail, err := team.Details(ctx)
		if err != nil {
			return nil, fmt.Errorf("getting details for team: %w", err)
		}
		details = append(details, detail)
	}
	return details, nil
}

func TeamsMenu(community domain.Community) content.View {
	return teamsMenuView{
		community: community,
	}
}

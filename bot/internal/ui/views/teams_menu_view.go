package views

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/conditions"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TeamsEmptyView() sources.Source[content.Content] {
	return sources.Static[content.Content](keyboards.Inline(
		content.Text("У вас пока нет команд. Создайте первую!"),
		buttons.Matrix(
			buttons.Row(buttons.CallbackButton("Create team", protocol.StartCreateTeamForm())),
			buttons.Row(buttons.CallbackButton("Main menu", protocol.ToMainMenu())),
		),
	))
}

type teamsView struct {
	teams sources.Source[domain.Teams]
}

func (t teamsView) Value(ctx context.Context, update telegram.Update) (content.Content, error) {
	teamsSource, err := t.teams.Value(ctx, update)
	if err != nil {
		return nil, fmt.Errorf("getting teams source: %w", err)
	}
	teams, err := teamsSource.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting teams details: %w", err)
	}
	return keyboards.Inline(
		content.Text("Select a team:"),
		t.teamsMatrix(teams).
			WithRow(buttons.Row(buttons.CallbackButton("Create team", protocol.StartCreateTeamForm()))).
			WithRow(buttons.Row(buttons.CallbackButton("Main menu", protocol.ToMainMenu()))),
	), nil
}

func (t teamsView) teamsMatrix(teams []domain.Team) buttons.ButtonMatrix[buttons.InlineButton] {
	matrix := buttons.Matrix[buttons.InlineButton]()
	for _, team := range teams {
		matrix = matrix.WithRow(buttons.Row(t.teamButton(team)))
	}
	return matrix
}

func (t teamsView) teamButton(team domain.Team) buttons.InlineButton {
	return buttons.CallbackButton(team.Name(), protocol.ToTeamMenu(team.ID()))
}

func TeamsView(teams sources.Source[domain.Teams]) sources.Source[content.Content] {
	return teamsView{teams: teams}
}

func TeamsMenu(comm domain.Community) sources.Source[content.Content] {
	return sources.IfElse(
		conditions.IsTeamsEmpty(comm),
		TeamsEmptyView(),
		TeamsView(CurrentTeams(comm)),
	)
}

type currentTeams struct {
	community domain.Community
}

func (c currentTeams) Value(ctx context.Context, update telegram.Update) (domain.Teams, error) {
	id, found := attributes.ChatID().Value(update)
	if !found {
		return nil, fmt.Errorf("chat ID not found in update")
	}
	user, err := c.community.UserWithTelegramID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting user with telegram ID %d: %w", id, err)
	}
	teams, err := user.CreatedTeams(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting created teams for user %d: %w", id, err)
	}
	return teams, nil
}

func CurrentTeams(comm domain.Community) sources.Source[domain.Teams] {
	return currentTeams{community: comm}
}

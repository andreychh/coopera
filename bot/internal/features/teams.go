package features

import (
	"context"
	"fmt"
	"iter"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/conditions"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/base"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
	"github.com/andreychh/coopera-bot/pkg/botlib/hsm"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TeamsEmptyView() sources.Source[content.Content] {
	return sources.Static[content.Content](keyboards.Inline(
		content.Text("У вас пока нет команд. Создайте первую!"),
		buttons.Matrix(
			buttons.Row(buttons.CallbackButton("Create team", protocol.StartForm("create_team"))),
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
	teamsDetails, err := teamsSource.Details(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting teams details: %w", err)
	}
	return keyboards.Inline(
		content.Text("Select a team:"),
		t.teamsMatrix(teamsDetails).
			WithRow(buttons.Row(buttons.CallbackButton("Create team", protocol.StartForm("create_team")))).
			WithRow(buttons.Row(buttons.CallbackButton("Main menu", protocol.ToMainMenu()))),
	), nil
}

func (t teamsView) teamsMatrix(teams iter.Seq2[int64, domain.TeamDetails]) buttons.ButtonMatrix[buttons.InlineButton] {
	matrix := buttons.Matrix[buttons.InlineButton]()
	for _, team := range teams {
		matrix = matrix.WithRow(buttons.Row(t.teamButton(team)))
	}
	return matrix
}

func (t teamsView) teamButton(team domain.TeamDetails) buttons.InlineButton {
	return buttons.CallbackButton(team.Name(), protocol.ToTeamMenu(team.ID()))
}

func TeamsView(teams sources.Source[domain.Teams]) sources.Source[content.Content] {
	return teamsView{teams: teams}
}

func TeamsMenu(comm domain.Community) sources.Source[content.Content] {
	return sources.IfElse(
		conditions.IsTeamsEmpty(comm),
		TeamsEmptyView(),
		TeamsView(domain.CurrentTeams(comm)),
	)
}

func TeamsMenuSpec(bot tg.Bot, c domain.Community) hsm.Spec {
	return hsm.Leaf(protocol.MenuTeams, hsm.CoreBehavior(
		base.EditOrSendContent(bot, TeamsMenu(c)),
		hsm.FirstHandled(
			hsm.JustIf(protocol.OnChangeMenu(protocol.MenuTeam), hsm.Transit(SpecTeamMenu)),
			hsm.JustIf(protocol.OnChangeMenu(protocol.MenuMain), hsm.Transit(SpecMainMenu)),
			hsm.JustIf(protocol.OnStartForm(protocol.FormCreateTeam), hsm.Transit(SpecCreateTeamForm)),
		),
		composition.Nothing(),
	))
}

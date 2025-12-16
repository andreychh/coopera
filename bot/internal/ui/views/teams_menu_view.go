package views

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/domain/conditions"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/formatting"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TeamsEmptyView() sources.Source[content.Content] {
	text := `👥 <b>Мои команды</b>

У вас пока нет ни одной команды.
Создайте свою первую команду, чтобы начать распределять задачи и добавлять участников!`

	return sources.Static[content.Content](keyboards.Inline(
		formatting.Formatted(content.Text(text), formatting.ParseModeHTML),
		buttons.Matrix(
			buttons.Row(buttons.CallbackButton("➕ Создать команду", protocol.StartCreateTeamForm())),
			buttons.Row(buttons.CallbackButton("🔙 Главное меню", protocol.ToMainMenu())),
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
	text := `👥 <b>Мои команды</b>

Выберите команду из списка ниже для управления задачами и участниками.`
	return keyboards.Inline(
		formatting.Formatted(content.Text(text), formatting.ParseModeHTML),
		t.teamsMatrix(teams).
			WithRow(buttons.Row(buttons.CallbackButton("➕ Создать команду", protocol.StartCreateTeamForm()))).
			WithRow(buttons.Row(buttons.CallbackButton("🔙 Главное меню", protocol.ToMainMenu()))),
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
	return buttons.CallbackButton(
		fmt.Sprintf("🏢 %s", team.Name()),
		protocol.ToTeamMenu(team.ID()),
	)
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
	user, exists, err := c.community.UserWithTelegramID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting user with telegram ID %d: %w", id, err)
	}
	if !exists {
		return nil, fmt.Errorf("user with telegram ID %d does not exist", id)
	}
	teams, err := user.Teams(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting created teams for user %d: %w", id, err)
	}
	return teams, nil
}

func CurrentTeams(comm domain.Community) sources.Source[domain.Teams] {
	return currentTeams{community: comm}
}

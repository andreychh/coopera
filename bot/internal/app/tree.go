package app

import (
	"log/slog"

	"github.com/andreychh/coopera-bot/internal/domain"
	domainactions "github.com/andreychh/coopera-bot/internal/domain/actions"
	uiactions "github.com/andreychh/coopera-bot/internal/ui/actions"
	"github.com/andreychh/coopera-bot/internal/ui/menu"
	"github.com/andreychh/coopera-bot/pkg/botlib/base"
	callbacksconditions "github.com/andreychh/coopera-bot/pkg/botlib/callbacks/conditions"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/dialogues"
	dialoguesactions "github.com/andreychh/coopera-bot/pkg/botlib/dialogues/actions"
	dialoguesconditions "github.com/andreychh/coopera-bot/pkg/botlib/dialogues/conditions"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms/actions"
	"github.com/andreychh/coopera-bot/pkg/botlib/logging"
	"github.com/andreychh/coopera-bot/pkg/botlib/routing"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	updatesconditions "github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
)

func Tree(bot tg.Bot, c domain.Community, d dialogues.Dialogues, f forms.Forms) core.Clause {
	return logging.LoggingClause(base.Recover(
		routing.If(updatesconditions.ChatTypeIs(updates.ChatTypePrivate),
			routing.FirstExecuted(
				routing.TerminalIf(
					composition.Not(dialoguesconditions.DialogueExists(d)),
					composition.Sequential(
						domainactions.CreateUser(c),
						dialoguesactions.StartNeutralDialog(d),
						uiactions.SendWelcomeMessage(bot),
						uiactions.SendMainMenu(bot),
					),
				),
				routing.TerminalIf(
					composition.All(
						updatesconditions.UpdateTypeIs(updates.UpdateTypeCallbackQuery),
						callbacksconditions.PrefixIs("change_menu"),
						callbacksconditions.ValueIs("menu_name", "main"),
					),
					uiactions.SendMainMenu(bot),
				),
				routing.TerminalIf(
					composition.All(
						updatesconditions.UpdateTypeIs(updates.UpdateTypeCallbackQuery),
						callbacksconditions.PrefixIs("change_menu"),
						callbacksconditions.ValueIs("menu_name", "teams"),
					),
					uiactions.SendTeamsMenu(bot, c),
				),
				createTeamForm(bot, c, d, f),
			),
		)),
		slog.Default(),
	)
}

func createTeamForm(bot tg.Bot, c domain.Community, d dialogues.Dialogues, f forms.Forms) core.Clause {
	return routing.FirstExecuted(
		routing.TerminalIf(
			composition.All(
				updatesconditions.UpdateTypeIs(updates.UpdateTypeCallbackQuery),
				dialoguesconditions.TopicIs(d, dialogues.TopicNeutral),
				callbacksconditions.PrefixIs("start_form"),
				callbacksconditions.ValueIs("form_name", "create_team"),
			),
			composition.Sequential(
				base.EditMessage(bot, menu.FormMenu()),
				base.SendContent(bot, content.Text("Please provide the name of your team.")),
				dialoguesactions.ChangeTopic(d, "form:create_team:name"),
			),
		),
		routing.TerminalIf(
			composition.All(
				updatesconditions.UpdateTypeIs(updates.UpdateTypeMessage),
				dialoguesconditions.TopicIs(d, "form:create_team:name"),
				composition.Not(updatesconditions.TextMatchesRegexp("^[A-Za-zА-Яа-я0-9_ -]{3,50}$")),
			),
			composition.Sequential(
				base.SendContent(bot, content.Text("Please provide the name of your team using 3 to 50 characters: letters, numbers, spaces, hyphens, or underscores.")),
			),
		),
		routing.TerminalIf(
			composition.All(
				updatesconditions.UpdateTypeIs(updates.UpdateTypeMessage),
				dialoguesconditions.TopicIs(d, "form:create_team:name"),
			),
			composition.Sequential(
				actions.SaveTextToField(f, "team_name"),
				domainactions.CreateTeam(f, c),
				base.SendContent(bot, content.Text("Great! Your team has been created.")),
				uiactions.SendMainMenu(bot),
				dialoguesactions.ChangeTopic(d, dialogues.TopicNeutral),
			),
		),
	)
}

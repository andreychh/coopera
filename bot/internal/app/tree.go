package app

import (
	"log/slog"

	"github.com/andreychh/coopera-bot/internal/domain"
	domainactions "github.com/andreychh/coopera-bot/internal/domain/actions"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/internal/ui/views"
	"github.com/andreychh/coopera-bot/pkg/botlib/base"
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
						domainactions.CreateUser(domain.IdempotencyCommunity(c)),
						dialoguesactions.StartNeutralDialog(d),
						base.SendContent(bot, views.WelcomeMessage()),
						base.SendContent(bot, views.MainMenu()),
					),
				),
				routing.TerminalIf(
					composition.All(
						dialoguesconditions.TopicIs(d, dialogues.TopicNeutral),
						composition.Any(
							composition.All(
								updatesconditions.UpdateTypeIs(updates.UpdateTypeMessage),
								updatesconditions.CommandIs("start"),
							),
							composition.All(
								updatesconditions.UpdateTypeIs(updates.UpdateTypeCallbackQuery),
								protocol.Navigation.On(protocol.MenuMain),
							),
						),
					),
					base.SendContent(bot, views.MainMenu()),
				),
				routing.TerminalIf(
					composition.All(
						updatesconditions.UpdateTypeIs(updates.UpdateTypeCallbackQuery),
						protocol.Navigation.On(protocol.MenuTeams),
					),
					base.SendContent(bot, views.TeamsMenu(c)),
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
				protocol.Forms.OnStart(protocol.FormCreateTeam),
			),
			composition.Sequential(
				base.EditMessage(bot, content.StaticView(content.Text("Please provide the name of your team."))),
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
				base.SendContent(
					bot,
					content.StaticView(
						content.Text("Please provide the name of your team using 3 to 50 characters: letters, numbers, spaces, hyphens, or underscores."),
					),
				),
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
				base.SendContent(
					bot,
					content.StaticView(
						content.Text("Great! Your team has been created."),
					),
				),
				base.SendContent(bot, views.TeamsMenu(c)),
				dialoguesactions.ChangeTopic(d, dialogues.TopicNeutral),
			),
		),
	)
}

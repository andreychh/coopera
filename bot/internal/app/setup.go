package app

import (
	"log/slog"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/ui"
	"github.com/andreychh/coopera-bot/pkg/botlib/base"
	"github.com/andreychh/coopera-bot/pkg/botlib/callbacks"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/dialogues"
	"github.com/andreychh/coopera-bot/pkg/botlib/engine"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
	"github.com/andreychh/coopera-bot/pkg/botlib/logging"
	"github.com/andreychh/coopera-bot/pkg/botlib/routing"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	"github.com/andreychh/coopera-bot/pkg/botlib/transport"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
)

func Store() keyvalue.Store {
	return keyvalue.MemoryStore()
}

func Dialogues(store keyvalue.Store) dialogues.Dialogues {
	return dialogues.KeyValueDialogues(store)
}

func Forms(store keyvalue.Store) forms.Forms {
	return forms.KeyValueForms(store)
}

func Client(token string) transport.Client {
	return transport.HTTPClient(token)
}

func Bot(client transport.Client) tg.Bot {
	return tg.NewBot(client)
}

func Community() domain.Community {
	return domain.MemoryCommunity{}
}

func Tree(bot tg.Bot, c domain.Community, d dialogues.Dialogues, f forms.Forms) core.Clause {
	return logging.LoggingClause(
		routing.FirstExecuted(
			routing.TerminalIf(
				composition.Not(dialogues.SafeDialogueExists(d)),
				composition.Sequential(
					domain.CreateUser(c),
					dialogues.StartNeutralDialog(d),
					domain.SendWelcomeMessage(bot),
					ui.SendMainMenu(bot),
				),
			),
			routing.TerminalIf(
				composition.All(
					callbacks.PrefixIs("change_menu"),
					callbacks.ParamIs("menu_name", "teams"),
				),
				ui.SendTeamsMenu(bot, c),
			),
			routing.TerminalIf(
				composition.All(
					callbacks.PrefixIs("change_menu"),
					callbacks.ParamIs("menu_name", "team"),
				),
				ui.SendTeamMenu(bot, c),
			),
			routing.TerminalIf(
				composition.All(
					dialogues.SafeTopicIs(d, dialogues.TopicNeutral),
					updates.SafeCommandIs("new_team"),
				),
				composition.Sequential(
					dialogues.ChangeTopic(d, "create_team-name"),
					base.SendContent(bot, content.Text("Let's create a new team! What is the name of your team?")),
				),
			),
			routing.TerminalIf(
				composition.All(
					dialogues.SafeTopicIs(d, "create_team-name"),
					composition.Not(updates.SafeTextMatchesRegexp("^[A-Za-zА-Яа-я0-9_ -]{3,50}$")),
				),
				composition.Sequential(
					base.SendContent(bot, content.Text(
						"The team name is invalid. Please provide a name between 3 and 50 characters, "+
							"using letters, numbers, spaces, hyphens, or underscores.",
					)),
				),
			),
			routing.TerminalIf(
				dialogues.SafeTopicIs(d, "create_team-name"),
				composition.Sequential(
					forms.SaveTextToField(f, "name"),
					dialogues.ChangeTopic(d, dialogues.TopicNeutral),
					base.SendContent(bot, content.Text("Great! Your team has been created.")),
				),
			),
		),
		slog.Default(),
	)
}

func Engine(token string, clause core.Clause) engine.Engine {
	return engine.ShutdownEngine(
		engine.SingleWorkerEngine(
			token, clause,
		),
	)
}

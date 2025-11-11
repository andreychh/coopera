package app

import (
	"fmt"
	"log/slog"

	"github.com/andreychh/coopera-bot/pkg/botlib/base"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/dialogues"
	"github.com/andreychh/coopera-bot/pkg/botlib/engine"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
	"github.com/andreychh/coopera-bot/pkg/botlib/logging"
	"github.com/andreychh/coopera-bot/pkg/botlib/routing"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func Bot(token string) (*telegram.BotAPI, error) {
	bot, err := telegram.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("initializing bot: %w", err)
	}
	bot.Debug = true
	return bot, nil
}

func Updates(bot *telegram.BotAPI) telegram.UpdatesChannel {
	u := telegram.NewUpdate(0)
	u.Timeout = 60
	return bot.GetUpdatesChan(u)
}

func Tree(bot *telegram.BotAPI, d dialogues.Dialogues, f forms.Forms) core.Clause {
	return logging.LoggingClause(
		routing.FirstExecuted(
			routing.TerminalIf(
				composition.All(
					composition.Not(dialogues.SafeDialogueExists(d)),
					updates.SafeCommandIs("start"),
				),
				composition.Sequential(
					dialogues.StartNeutralDialog(d),
					base.SendMessage(bot, "Hello, new user! To create a team use /new_team command."),
				),
			),
			routing.TerminalIf(
				composition.All(
					dialogues.SafeTopicIs(d, dialogues.TopicNeutral),
					updates.SafeCommandIs("new_team"),
				),
				composition.Sequential(
					dialogues.ChangeTopic(d, "create_team-name"),
					base.SendMessage(bot, "Let's create a new team! What is the name of your team?"),
				),
			),
			routing.TerminalIf(
				composition.All(
					dialogues.SafeTopicIs(d, "create_team-name"),
					composition.Not(updates.SafeTextMatchesRegexp("^[A-Za-zА-Яа-я0-9_ -]{3,50}$")),
				),
				composition.Sequential(
					base.SendMessage(bot,
						"The team name is invalid. Please provide a name between 3 and 50 characters, "+
							"using letters, numbers, spaces, hyphens, or underscores.",
					),
				),
			),
			routing.TerminalIf(
				dialogues.SafeTopicIs(d, "create_team-name"),
				composition.Sequential(
					forms.SaveTextToField(f, "name"),
					dialogues.ChangeTopic(d, dialogues.TopicNeutral),
					base.SendMessage(bot, "Great! Your team has been created."),
				),
			),
		),
		slog.Default(),
	)
}

func Engine(clause core.Clause, updates telegram.UpdatesChannel) engine.Engine {
	return engine.ShutdownEngine(
		engine.SingleWorkerEngine(
			clause, updates,
		),
	)
}

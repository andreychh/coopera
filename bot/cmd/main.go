package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/andreychh/coopera-bot/internal/app"
	"github.com/andreychh/coopera-bot/pkg/botlib/base"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/hsm"
	"github.com/andreychh/coopera-bot/pkg/botlib/logging"
	"github.com/andreychh/coopera-bot/pkg/botlib/routing"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/conditions"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot := app.Bot(token)
	community := app.Community(os.Getenv("BACKEND_API_URL"))
	store := app.Store()
	sessions := app.Sessions(store)
	forms := app.Forms(store)
	tree := app.Tree(bot, community, forms)
	graph := core.Must(hsm.NewCompiler(tree).Graph())
	engine := app.Engine(token,
		composition.Run(
			logging.LoggingClause(
				base.Recover(
					routing.If(
						composition.All(
							conditions.ChatTypeIs(updates.ChatTypePrivate),
							composition.Any(
								conditions.UpdateTypeIs(updates.UpdateTypeMessage),
								conditions.UpdateTypeIs(updates.UpdateTypeCallbackQuery),
							),
						),
						hsm.NewEngine(sessions, graph, bot),
					),
				),
				slog.Default(),
			),
		),
	)
	engine.Start(context.Background())
}

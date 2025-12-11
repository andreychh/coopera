package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/andreychh/coopera-bot/internal/app"
	"github.com/andreychh/coopera-bot/pkg/botlib/base"
	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/hsm"
	"github.com/andreychh/coopera-bot/pkg/botlib/logging"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot := app.Bot(token)
	community := app.Community(os.Getenv("BACKEND_API_URL"))
	store := app.Store()
	sessions := app.Sessions(store)
	forms := app.Forms(store)
	tree := app.Tree(bot, community, forms)
	graph, err := hsm.NewCompiler(tree).Graph()
	if err != nil {
		panic(err)
	}
	engine := app.Engine(token,
		composition.Run(
			base.Recover(
				logging.LoggingClause(
					hsm.NewEngine(sessions, graph),
					slog.Default(),
				),
			),
		),
	)
	engine.Start(context.Background())
}

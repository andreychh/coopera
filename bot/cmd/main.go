package main

import (
	"context"
	"os"

	"github.com/andreychh/coopera-bot/internal/app"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot := app.Bot(app.Client(token))
	community := app.Community(os.Getenv("BACKEND_API_URL"))
	store := app.Store()
	dialogues := app.Dialogues(store)
	forms := app.Forms(store)
	tree := app.Tree(bot, community, dialogues, forms)
	engine := app.Engine(token, tree)
	engine.Start(context.Background())
}

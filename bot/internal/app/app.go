package app

import (
	"context"
	"net/http"
	"os"

	"github.com/andreychh/coopera-bot/internal/transport"
)

func Start() error {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	client := Client(token)
	bot := Bot(client)
	community := Community(transport.HTTPClient(os.Getenv("BACKEND_API_URL"), &http.Client{}))
	store := Store()
	dialogues := Dialogues(store)
	forms := Forms(store)
	tree := Tree(bot, community, dialogues, forms)
	engine := Engine(token, tree)
	engine.Start(context.Background())
	return nil
}

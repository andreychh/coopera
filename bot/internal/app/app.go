package app

import (
	"context"
	"os"
)

func Start() error {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	client := Client(token)
	bot := Bot(client)
	community := Community()
	store := Store()
	dialogues := Dialogues(store)
	forms := Forms(store)
	tree := Tree(bot, community, dialogues, forms)
	engine := Engine(token, tree)
	engine.Start(context.Background())
	return nil
}

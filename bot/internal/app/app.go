package app

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Start() error {
	_ = godotenv.Load("../.env.bot")
	bot, err := Bot(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		return fmt.Errorf("initializing telegram bot: %w", err)
	}
	store := Store()
	dialogues := Dialogues(store)
	tree := Tree(bot, dialogues)
	updates := Updates(bot)
	engine := Engine(tree, updates)
	engine.Start(context.Background())
	return nil
}

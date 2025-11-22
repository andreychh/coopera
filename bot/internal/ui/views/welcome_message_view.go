package views

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type welcomeMessageView struct{}

func (w welcomeMessageView) Render(ctx context.Context, update telegram.Update) (content.Content, error) {
	return content.Text("Welcome to Coopera Bot! Use the menu below to navigate."), nil
}

func WelcomeMessage() content.View {
	return welcomeMessageView{}
}

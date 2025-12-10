package views

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
)

func WelcomeMessage() sources.Source[content.Content] {
	return sources.Static(content.Text("Welcome to Coopera Bot! Use the menu below to navigate."))
}

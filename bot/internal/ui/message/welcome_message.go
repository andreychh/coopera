package message

import "github.com/andreychh/coopera-bot/pkg/botlib/content"

func WelcomeMessage() content.Content {
	return content.Text("Welcome to the community! We're glad to have you here. Feel free to explore and engage with others.")
}

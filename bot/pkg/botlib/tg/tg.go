package tg

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/content"
)

type Bot interface {
	Chat(id int64) Chat
	AnswerCallbackQuery(ctx context.Context, id string) error
}

type Chat interface {
	Send(ctx context.Context, cnt content.Content) error
	Message(id int) Message
}

type Message interface {
	Edit(ctx context.Context, content content.Content) error
	Delete(ctx context.Context) error
}

package hsm

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Behavior interface {
	Prompt(ctx context.Context, u telegram.Update) error
	React(ctx context.Context, u telegram.Update) (Decision, error)
	Cleanup(ctx context.Context, u telegram.Update) error
}

type Compiler interface {
	Graph() (Graph, error)
}

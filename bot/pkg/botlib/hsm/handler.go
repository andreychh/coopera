package hsm

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler interface {
	Handle(ctx context.Context, u telegram.Update) (Decision, error)
}

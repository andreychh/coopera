package sources

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Source[T any] interface {
	Value(ctx context.Context, update telegram.Update) (T, error)
}

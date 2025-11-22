package content

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type View interface {
	Render(ctx context.Context, update telegram.Update) (Content, error)
}

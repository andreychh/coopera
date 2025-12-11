package hsm

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type firstHandledHandler struct {
	handlers []Handler
}

func (f firstHandledHandler) Handle(ctx context.Context, u telegram.Update) (Decision, error) {
	for _, handler := range f.handlers {
		decision, err := handler.Handle(ctx, u)
		if err != nil {
			return nil, err
		}
		if decision.Handled() {
			return decision, nil
		}
	}
	return Pass(), nil
}

func FirstHandled(handlers ...Handler) Handler {
	return firstHandledHandler{
		handlers: handlers,
	}
}

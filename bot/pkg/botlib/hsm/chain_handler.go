package hsm

import (
	"context"
	"iter"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type chainHandler struct {
	iter iter.Seq[State]
}

func (c chainHandler) Handle(ctx context.Context, u telegram.Update) (Decision, error) {
	for handler := range c.iter {
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

func Chain(iter iter.Seq[State]) Handler {
	return chainHandler{
		iter: iter,
	}
}

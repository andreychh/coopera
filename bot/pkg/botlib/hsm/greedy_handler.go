package hsm

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type greedyHandler struct {
	origin Handler
}

func (g greedyHandler) Handle(ctx context.Context, u telegram.Update) (Decision, error) {
	decision, err := g.origin.Handle(ctx, u)
	if err != nil {
		return nil, err
	}
	if decision.Handled() {
		return decision, nil
	}
	return Stay(), nil
}

func Greedy(origin Handler) Handler {
	return greedyHandler{
		origin: origin,
	}
}

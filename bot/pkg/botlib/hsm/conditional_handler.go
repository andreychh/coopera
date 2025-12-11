package hsm

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type conditionalHandler struct {
	condition core.Condition
	handler   Handler
}

func (c conditionalHandler) Handle(ctx context.Context, u telegram.Update) (Decision, error) {
	held, err := c.condition.Holds(ctx, u)
	if err != nil {
		return nil, err
	}
	if !held {
		return Pass(), nil
	}
	return c.handler.Handle(ctx, u)
}

func If(condition core.Condition, handler Handler) Handler {
	return conditionalHandler{
		condition: condition,
		handler:   handler,
	}
}

func JustIf(condition core.Condition, decision Decision) Handler {
	return If(condition, Just(decision))
}

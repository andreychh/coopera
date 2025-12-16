package sources

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type conditionalSource[T any] struct {
	condition core.Condition
	true      Source[T]
	false     Source[T]
}

func (c conditionalSource[T]) Value(ctx context.Context, update telegram.Update) (T, error) {
	holds, err := c.condition.Holds(ctx, update)
	if err != nil {
		return core.Zero[T](), err
	}
	if holds {
		return c.true.Value(ctx, update)
	}
	return c.false.Value(ctx, update)
}

func IfElse[T any](condition core.Condition, true Source[T], false Source[T]) Source[T] {
	return conditionalSource[T]{
		condition: condition,
		true:      true,
		false:     false,
	}
}

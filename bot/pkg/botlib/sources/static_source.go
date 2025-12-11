package sources

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type staticSource[T any] struct {
	value T
}

func (s staticSource[T]) Value(_ context.Context, _ telegram.Update) (T, error) {
	return s.value, nil
}

func Static[T any](value T) Source[T] {
	return staticSource[T]{
		value: value,
	}
}

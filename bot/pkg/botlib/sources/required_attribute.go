package sources

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type requiredAttribute[T any] struct {
	attribute attributes.Attribute[T]
}

func (r requiredAttribute[T]) Value(_ context.Context, update telegram.Update) (T, error) {
	value, found := r.attribute.Value(update)
	if !found {
		return core.Zero[T](), fmt.Errorf("attribute %T not found", r.attribute)
	}
	return value, nil
}

func Required[T any](attribute attributes.Attribute[T]) Source[T] {
	return requiredAttribute[T]{
		attribute: attribute,
	}
}

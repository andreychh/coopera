package attributes

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Attribute[T any] interface {
	Value(update telegram.Update) (value T, found bool)
}

type firstFoundAttribute[T any] struct {
	attributes []Attribute[T]
}

func (f firstFoundAttribute[T]) Value(update telegram.Update) (T, bool) {
	for _, attr := range f.attributes {
		value, found := attr.Value(update)
		if found {
			return value, true
		}
	}
	return core.Zero[T](), false
}

func OneOf[T any](attributes ...Attribute[T]) Attribute[T] {
	return firstFoundAttribute[T]{
		attributes: attributes,
	}
}

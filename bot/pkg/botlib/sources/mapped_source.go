package sources

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type mappedSource[I, O any] struct {
	origin  Source[I]
	mapping Mapping[I, O]
}

func (m mappedSource[I, O]) Value(ctx context.Context, update telegram.Update) (O, error) {
	origin, err := m.origin.Value(ctx, update)
	if err != nil {
		return core.Zero[O](), fmt.Errorf("getting origin (%T): %w", m.origin, err)
	}
	mapped, err := m.mapping.Apply(ctx, origin)
	if err != nil {
		return core.Zero[O](), fmt.Errorf("mapping: %w", err)
	}
	return mapped, nil
}

func Mapped[I, O any](source Source[I], mapping Mapping[I, O]) Source[O] {
	return mappedSource[I, O]{
		origin:  source,
		mapping: mapping,
	}
}

func Map[I, O any](source Source[I], mapping func(context.Context, I) (O, error)) Source[O] {
	return Mapped(source, MappingFunc(
		func(ctx context.Context, input I) (O, error) {
			return mapping(ctx, input)
		},
	))
}

func ErrMap[I, O any](source Source[I], mapping func(I) (O, error)) Source[O] {
	return Mapped(source, MappingFunc(
		func(_ context.Context, input I) (O, error) {
			return mapping(input)
		},
	))
}

func PureMap[I, O any](source Source[I], mapping func(I) O) Source[O] {
	return Mapped(source, MappingFunc(
		func(_ context.Context, input I) (O, error) {
			return mapping(input), nil
		},
	))
}

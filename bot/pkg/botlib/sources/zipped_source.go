package sources

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type zippedSource[A, B, O any] struct {
	a       Source[A]
	b       Source[B]
	zipping Zipping[A, B, O]
}

func (z zippedSource[A, B, O]) Value(ctx context.Context, update telegram.Update) (O, error) {
	a, err := z.a.Value(ctx, update)
	if err != nil {
		return core.Zero[O](), fmt.Errorf("getting 1st zip arg (%T): %w", z.a, err)
	}
	b, err := z.b.Value(ctx, update)
	if err != nil {
		return core.Zero[O](), fmt.Errorf("getting 2nd zip arg (%T): %w", z.b, err)
	}
	zipped, err := z.zipping.Apply(ctx, a, b)
	if err != nil {
		return core.Zero[O](), fmt.Errorf("zipping: %w", err)
	}
	return zipped, nil
}

func Zipped[A, B, O any](a Source[A], b Source[B], zipping Zipping[A, B, O]) Source[O] {
	return zippedSource[A, B, O]{
		a:       a,
		b:       b,
		zipping: zipping,
	}
}

func Zip[A, B, O any](a Source[A], b Source[B], zipping func(context.Context, A, B) (O, error)) Source[O] {
	return Zipped(a, b, ZippingFunc(
		func(ctx context.Context, a A, b B) (O, error) {
			return zipping(ctx, a, b)
		},
	))
}

func ErrZip[A, B, O any](a Source[A], b Source[B], zipping func(A, B) (O, error)) Source[O] {
	return Zipped(a, b, ZippingFunc(
		func(_ context.Context, a A, b B) (O, error) {
			return zipping(a, b)
		},
	))
}

func PureZip[A, B, O any](a Source[A], b Source[B], zipping func(A, B) O) Source[O] {
	return Zipped(a, b, ZippingFunc(
		func(_ context.Context, a A, b B) (O, error) {
			return zipping(a, b), nil
		},
	))
}

package sources

import (
	"context"
)

type Zipping[A, B, O any] interface {
	Apply(ctx context.Context, a A, b B) (output O, err error)
}

type zippingFunc[A, B, O any] func(ctx context.Context, a A, b B) (O, error)

func (f zippingFunc[A, B, O]) Apply(ctx context.Context, a A, b B) (O, error) {
	return f(ctx, a, b)
}

func ZippingFunc[A, B, O any](fn func(context.Context, A, B) (O, error)) Zipping[A, B, O] {
	return zippingFunc[A, B, O](fn)
}

package sources

import (
	"context"
)

type Application[A, B any] interface {
	Apply(ctx context.Context, a A, b B) error
}
type applicationFunc[A, B any] func(context.Context, A, B) error

func (f applicationFunc[A, B]) Apply(ctx context.Context, a A, b B) error {
	return f(ctx, a, b)
}

func ApplicationFunc[A, B any](fn func(context.Context, A, B) error) Application[A, B] {
	return applicationFunc[A, B](fn)
}

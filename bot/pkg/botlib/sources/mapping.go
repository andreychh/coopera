package sources

import "context"

type Mapping[I, O any] interface {
	Apply(ctx context.Context, input I) (output O, err error)
}
type mappingFunc[I, O any] func(ctx context.Context, input I) (O, error)

func (f mappingFunc[I, O]) Apply(ctx context.Context, input I) (O, error) {
	return f(ctx, input)
}

func MappingFunc[I, O any](fn func(context.Context, I) (O, error)) Mapping[I, O] {
	return mappingFunc[I, O](fn)
}

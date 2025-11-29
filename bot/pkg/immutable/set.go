package immutable

import (
	"iter"
)

type set[T comparable] struct {
	items Map[T, struct{}]
}

func (s set[T]) Has(item T) bool {
	_, exists := s.items.Get(item)
	return exists
}

func (s set[T]) Add(item T) Set[T] {
	return set[T]{items: s.items.With(item, struct{}{})}
}

func (s set[T]) Remove(item T) Set[T] {
	return set[T]{items: s.items.Remove(item)}
}

func (s set[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		s.items.All()(func(key T, _ struct{}) bool {
			return yield(key)
		})
	}
}

func (s set[T]) Len() int {
	return s.items.Len()
}

func EmptySet[T comparable]() Set[T] {
	return set[T]{items: EmptyMap[T, struct{}]()}
}

func SetOf[T comparable](origin ...T) Set[T] {
	copied := make(map[T]struct{}, len(origin))
	for _, item := range origin {
		copied[item] = struct{}{}
	}
	return set[T]{items: MapOf(copied)}
}

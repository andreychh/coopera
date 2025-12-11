package immutable

import (
	"iter"
)

type reversedSlice[T any] struct {
	origin Slice[T]
}

func (r reversedSlice[T]) Get(index int) T {
	return r.origin.Get(r.origin.Len() - 1 - index)
}

func (r reversedSlice[T]) Insert(index int, item T) Slice[T] {
	return ReversedSlice(r.origin.Insert(r.origin.Len()-index, item))
}

func (r reversedSlice[T]) Replace(index int, item T) Slice[T] {
	return ReversedSlice(r.origin.Replace(r.origin.Len()-1-index, item))
}

func (r reversedSlice[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := r.origin.Len() - 1; i >= 0; i-- {
			if !yield(r.origin.Len()-1-i, r.origin.Get(i)) {
				return
			}
		}
	}
}

func (r reversedSlice[T]) Len() int {
	return r.origin.Len()
}

func ReversedSlice[T any](origin Slice[T]) Slice[T] {
	return reversedSlice[T]{origin: origin}
}

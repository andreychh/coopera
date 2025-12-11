package immutable

import (
	"iter"
	"slices"
)

type immutableSlice[T any] struct {
	items []T
}

func (i immutableSlice[T]) Get(index int) T {
	return i.items[index]
}

func (i immutableSlice[T]) Insert(index int, item T) Slice[T] {
	copied := make([]T, len(i.items)+1)
	copy(copied, i.items[:index])
	copied[index] = item
	copy(copied[index+1:], i.items[index:])
	return immutableSlice[T]{items: copied}
}

func (i immutableSlice[T]) Replace(index int, item T) Slice[T] {
	copied := make([]T, len(i.items))
	copy(copied, i.items)
	copied[index] = item
	return immutableSlice[T]{items: copied}
}

func (i immutableSlice[T]) All() iter.Seq2[int, T] {
	return slices.All(i.items)
}

func (i immutableSlice[T]) Len() int {
	return len(i.items)
}

func EmptySlice[T any]() Slice[T] {
	return immutableSlice[T]{items: make([]T, 0)}
}

func SliceOf[T any](origin []T) Slice[T] {
	copied := make([]T, len(origin))
	copy(copied, origin)
	return immutableSlice[T]{items: copied}
}

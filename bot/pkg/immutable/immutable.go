package immutable

import "iter"

type Map[K comparable, V any] interface {
	Get(key K) (V, bool)
	With(key K, value V) Map[K, V]
	Remove(key K) Map[K, V]
	All() iter.Seq2[K, V]
	Len() int
}

type Slice[T any] interface {
	Get(index int) T
	Insert(index int, item T) Slice[T]
	Replace(index int, item T) Slice[T]
	All() iter.Seq2[int, T]
	Len() int
}

type Set[T comparable] interface {
	Has(item T) bool
	Add(item T) Set[T]
	Remove(item T) Set[T]
	All() iter.Seq[T]
	Len() int
}

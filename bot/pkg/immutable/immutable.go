package immutable

import "iter"

type Map[K comparable, V any] interface {
	Get(key K) (V, bool)
	With(key K, value V) Map[K, V]
	Remove(key K) Map[K, V]
	All() iter.Seq2[K, V]
	Len() int
}

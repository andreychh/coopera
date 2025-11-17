package slices

// With appends elements to a slice and returns a new slice.
// todo: replace with ImmutableSlice[T any] type
func With[T any](slice []T, elements ...T) []T {
	result := make([]T, len(slice), len(slice)+len(elements))
	copy(result, slice)
	result = append(result, elements...)
	return result
}

func WithReplaced[T any](slice []T, index int, value T) []T {
	if index < 0 || index >= len(slice) {
		return With(slice)
	}
	result := make([]T, len(slice))
	copy(result, slice)
	result[index] = value
	return result
}

package slices

func With[T any](slice []T, elements ...T) []T {
	result := make([]T, 0, len(slice)+len(elements))
	result = append(result, slice...)
	result = append(result, elements...)
	return result
}

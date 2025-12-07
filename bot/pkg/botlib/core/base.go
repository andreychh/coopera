package core

func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func Zero[T any]() T {
	var z T
	return z
}

func Ptr[T any](v T) *T {
	return &v
}

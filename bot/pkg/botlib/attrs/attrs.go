package attrs

type Attribute[T any] interface {
	Value() (value T, exists bool)
}

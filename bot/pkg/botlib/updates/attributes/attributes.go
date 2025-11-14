package attributes

import "errors"

var ErrAttributeNotFound = errors.New("attribute not found")

type Attribute[T any] interface {
	Exists() bool
	Value() (T, error)
}

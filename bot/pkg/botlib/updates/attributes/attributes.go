package attributes

import "errors"

type Attribute[T any] interface {
	Exists() bool
	Value() (T, error)
}

var ErrAttributeNotFound = errors.New("attribute not found")

package repr

import "errors"

type Encodable interface {
	Encode() ([]byte, error)
	Update(path Path, value Encodable) (Encodable, error)
}

type Array interface {
	Encodable
	WithElement(element Encodable) Array
	Extend(array Array) Array
	AsSlice() []Encodable
}

type Object interface {
	Encodable
	WithField(key string, value Encodable) Object
	Merge(other Object) Object
	AsMap() map[string]Encodable
}

type Path interface {
	Empty() bool
	Index() (int, error)
	Key() (string, error)
	Tail() Path
}

var ErrCannotUpdate = errors.New("cannot update into this type")

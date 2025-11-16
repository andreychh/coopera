package repr

import "errors"

type Primitive interface {
	Marshal() ([]byte, error)
}

type Structure interface {
	Marshal() ([]byte, error)
	At(path Path) (Structure, error)
	Update(path Path, value Structure) (Structure, error)
	Extend(path Path, other Structure) (Structure, error)
}

type Path interface {
	Empty() bool
	Head() (Segment, error)
	Tail() Path
}

type Segment interface {
	Index() (int, bool)
	Key() (string, bool)
}

var ErrCannotGet = errors.New("cannot get value at the specified path")
var ErrCannotUpdate = errors.New("cannot update value at the specified path")
var ErrCannotMerge = errors.New("cannot merge value at the specified path")

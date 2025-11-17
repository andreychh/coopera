package json

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
)

type structure struct {
	origin repr.Primitive
}

func (s structure) Marshal() ([]byte, error) {
	return s.origin.Marshal()
}

func (s structure) At(path repr.Path) (repr.Structure, error) {
	if !path.Empty() {
		return nil, repr.ErrCannotGet
	}
	return s, nil
}

func (s structure) Update(path repr.Path, value repr.Structure) (repr.Structure, error) {
	if !path.Empty() {
		return nil, repr.ErrCannotUpdate
	}
	return value, nil
}

func (s structure) Extend(_ repr.Path, _ repr.Structure) (repr.Structure, error) {
	return nil, repr.ErrCannotMerge
}

func StructureOf(primitive repr.Primitive) repr.Structure {
	return structure{origin: primitive}
}

package json

import (
	"strconv"

	"github.com/andreychh/coopera-bot/pkg/repr"
)

type boolean struct {
	value bool
}

func (b boolean) Marshal() ([]byte, error) {
	return []byte(strconv.FormatBool(b.value)), nil
}

func Boolean(value bool) repr.Primitive {
	return boolean{value: value}
}

func Bool(value bool) repr.Structure {
	return StructureOf(Boolean(value))
}

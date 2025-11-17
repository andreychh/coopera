package json

import (
	"strconv"

	"github.com/andreychh/coopera-bot/pkg/repr"
)

type integer struct {
	value int64
}

func (f integer) Marshal() ([]byte, error) {
	return []byte(strconv.FormatInt(f.value, 10)), nil
}

func Integer(value int64) repr.Primitive {
	return integer{value: value}
}

func Int(value int) repr.Structure {
	return StructureOf(Integer(int64(value)))
}

func I64(value int64) repr.Structure {
	return StructureOf(Integer(value))
}

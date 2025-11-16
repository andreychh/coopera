package json

import (
	"strconv"

	"github.com/andreychh/coopera-bot/pkg/repr"
)

type float struct {
	value float64
}

func (f float) Marshal() ([]byte, error) {
	return []byte(strconv.FormatFloat(f.value, 'f', -1, 64)), nil
}

func Float(value float64) repr.Primitive {
	return float{value: value}
}

func F64(value float64) repr.Structure {
	return StructureOf(Float(value))
}

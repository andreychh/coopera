package json

import (
	"strconv"

	"github.com/andreychh/coopera-bot/pkg/repr"
)

type string_ struct {
	value string
}

func (s string_) Marshal() ([]byte, error) {
	return []byte(strconv.Quote(s.value)), nil
}

func String(value string) repr.Primitive {
	return string_{value: value}
}

func Str(value string) repr.Structure {
	return StructureOf(String(value))
}

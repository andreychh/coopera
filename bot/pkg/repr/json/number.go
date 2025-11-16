package json

import (
	"strconv"

	"github.com/andreychh/coopera-bot/pkg/repr"
)

type number struct {
	value float64
}

func (n number) Marshal() ([]byte, error) {
	return []byte(strconv.FormatFloat(n.value, 'f', -1, 64)), nil
}
func Number(value float64) repr.Primitive {
	return number{value: value}
}

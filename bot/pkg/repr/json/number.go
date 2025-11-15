package json

import (
	"strconv"

	"github.com/andreychh/coopera-bot/pkg/repr"
)

type number struct {
	value float64
}

func (n number) Encode() ([]byte, error) {
	return []byte(strconv.FormatFloat(n.value, 'f', -1, 64)), nil
}

func Number(value float64) repr.Encodable {
	return number{value: value}
}

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

func (n number) Update(path repr.Path, value repr.Encodable) (repr.Encodable, error) {
	if !path.Empty() {
		return nil, repr.ErrCannotUpdate
	}
	return value, nil
}

func Number(value float64) repr.Encodable {
	return number{value: value}
}

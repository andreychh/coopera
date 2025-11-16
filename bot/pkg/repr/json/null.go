package json

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
)

type null struct{}

func (n null) Encode() ([]byte, error) {
	return []byte("null"), nil
}

func (n null) Update(path repr.Path, value repr.Encodable) (repr.Encodable, error) {
	if !path.Empty() {
		return nil, repr.ErrCannotUpdate
	}
	return value, nil
}

func Null() repr.Encodable {
	return null{}
}

package json

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
)

type null struct{}

func (n null) Encode() ([]byte, error) {
	return []byte("null"), nil
}

func Null() repr.Encodable {
	return null{}
}

package json

import (
	"encoding/json"

	"github.com/andreychh/coopera-bot/pkg/repr"
)

type string_ struct {
	value string
}

func (s string_) Encode() ([]byte, error) {
	return json.Marshal(s.value)
}

func (s string_) Update(path repr.Path, value repr.Encodable) (repr.Encodable, error) {
	if !path.Empty() {
		return nil, repr.ErrCannotUpdate
	}
	return value, nil
}

func String(value string) repr.Encodable {
	return string_{value: value}
}

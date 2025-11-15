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

func String(value string) repr.Encodable {
	return string_{value: value}
}

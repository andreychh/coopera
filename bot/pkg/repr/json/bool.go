package json

import (
	"strconv"

	"github.com/andreychh/coopera-bot/pkg/repr"
)

type boolean struct {
	value bool
}

func (b boolean) Encode() ([]byte, error) {
	return []byte(strconv.FormatBool(b.value)), nil
}

func (b boolean) Update(path repr.Path, value repr.Encodable) (repr.Encodable, error) {
	if !path.Empty() {
		return nil, repr.ErrCannotUpdate
	}
	return value, nil
}

func Boolean(value bool) repr.Encodable {
	return boolean{value: value}
}

package json

import (
	"bytes"

	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/slices"
)

type array struct {
	elements []repr.Encodable
}

func (a array) Encode() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, elem := range a.elements {
		if i > 0 {
			buf.WriteByte(',')
		}
		encoded, err := elem.Encode()
		if err != nil {
			return nil, err
		}
		buf.Write(encoded)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func (a array) WithElement(element repr.Encodable) repr.Array {
	return array{elements: slices.With(a.elements, element)}
}

func (a array) Extend(other repr.Array) repr.Array {
	return array{elements: slices.With(a.elements, other.AsSlice()...)}
}

func (a array) AsSlice() []repr.Encodable {
	return slices.With(a.elements)
}

func Array(elements ...repr.Encodable) repr.Array {
	return array{elements: slices.With(elements)}
}

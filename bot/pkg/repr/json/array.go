package json

import (
	"bytes"
	"fmt"

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

func (a array) Update(path repr.Path, value repr.Encodable) (repr.Encodable, error) {
	if path.Empty() {
		return value, nil
	}
	index, err := path.Index()
	if err != nil {
		return nil, err
	}
	if !a.correctIndex(index) {
		return nil, fmt.Errorf("index %d out of bounds", index)
	}
	updated, err := a.elements[index].Update(path.Tail(), value)
	if err != nil {
		return nil, err
	}
	return array{elements: slices.WithReplaced(a.elements, index, updated)}, nil
}

func (a array) correctIndex(index int) bool {
	return index >= 0 && index < len(a.elements)
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

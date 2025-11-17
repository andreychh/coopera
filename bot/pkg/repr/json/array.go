package json

import (
	"bytes"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/slices"
)

type array struct {
	elements []repr.Structure
}

func (a array) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, elem := range a.elements {
		if i > 0 {
			buf.WriteByte(',')
		}
		data, err := elem.Marshal()
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func (a array) At(path repr.Path) (repr.Structure, error) {
	if path.Empty() {
		return a, nil
	}
	index, err := a.index(path)
	if err != nil {
		return nil, err
	}
	return a.elements[index].At(path.Tail())
}

func (a array) Update(path repr.Path, value repr.Structure) (repr.Structure, error) {
	if path.Empty() {
		return value, nil
	}
	index, err := a.index(path)
	if err != nil {
		return nil, err
	}
	updated, err := a.elements[index].Update(path.Tail(), value)
	if err != nil {
		return nil, err
	}
	return array{elements: slices.WithReplaced(a.elements, index, updated)}, nil
}

func (a array) Extend(path repr.Path, other repr.Structure) (repr.Structure, error) {
	if path.Empty() {
		otherArray, ok := other.(array)
		if !ok {
			return nil, fmt.Errorf("can only extend array with array, got %T", other)
		}
		return array{elements: slices.With(a.elements, otherArray.elements...)}, nil
	}
	index, err := a.index(path)
	if err != nil {
		return nil, err
	}
	extended, err := a.elements[index].Extend(path.Tail(), other)
	if err != nil {
		return nil, err
	}
	return array{elements: slices.WithReplaced(a.elements, index, extended)}, nil
}

func (a array) index(path repr.Path) (int, error) {
	head, err := path.Head()
	if err != nil {
		return 0, fmt.Errorf("getting path head: %w", err)
	}
	index, ok := head.Index()
	if !ok {
		return 0, fmt.Errorf("array path must be index, got %v", head)
	}
	if index < 0 || index >= len(a.elements) {
		return 0, fmt.Errorf("index %d out of bounds", index)
	}
	return index, nil
}

func Array(elements ...repr.Structure) repr.Structure {
	return array{elements: slices.With(elements)}
}

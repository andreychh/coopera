package json

import (
	"bytes"

	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/slices"
)

type object struct {
	fields []field
}

func (o object) Encode() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, fld := range o.fields {
		if i > 0 {
			buf.WriteByte(',')
		}
		encoded, err := fld.Encode()
		if err != nil {
			return nil, err
		}
		buf.Write(encoded)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (o object) Update(path repr.Path, value repr.Encodable) (repr.Encodable, error) {
	if path.Empty() {
		return value, nil
	}
	key, err := path.Key()
	if err != nil {
		return nil, err
	}
	for i, fld := range o.fields {
		if fld.key == key {
			updated, err := fld.value.Update(path.Tail(), value)
			if err != nil {
				return nil, err
			}
			return object{fields: slices.WithReplaced(o.fields, i, field{key: key, value: updated})}, nil
		}
	}
	return nil, repr.ErrCannotUpdate
}

func (o object) WithField(key string, value repr.Encodable) repr.Object {
	return object{fields: slices.With(o.fields, field{key: key, value: value})}
}

func (o object) Merge(other repr.Object) repr.Object {
	otherMap := other.AsMap()
	merged := make([]field, 0, len(o.fields)+len(otherMap))
	seen := make(map[string]bool)
	for _, f := range o.fields {
		if newValue, exists := otherMap[f.key]; exists {
			merged = append(merged, field{key: f.key, value: newValue})
		} else {
			merged = append(merged, f)
		}
		seen[f.key] = true
	}
	for key, value := range otherMap {
		if !seen[key] {
			merged = append(merged, field{key: key, value: value})
		}
	}
	return object{fields: merged}
}

func (o object) AsMap() map[string]repr.Encodable {
	result := make(map[string]repr.Encodable, len(o.fields))
	for _, f := range o.fields {
		result[f.key] = f.value
	}
	return result
}

func Object() repr.Object {
	return object{fields: []field{}}
}

type field struct {
	key   string
	value repr.Encodable
}

func (f field) Update(path repr.Path, value repr.Encodable) repr.Encodable {
	// TODO implement me
	panic("implement me")
}

func (f field) Encode() ([]byte, error) {
	key, err := String(f.key).Encode()
	if err != nil {
		return nil, err
	}
	value, err := f.value.Encode()
	if err != nil {
		return nil, err
	}
	result := make([]byte, 0, len(key)+1+len(value))
	result = append(result, key...)
	result = append(result, ':')
	result = append(result, value...)
	return result, nil
}

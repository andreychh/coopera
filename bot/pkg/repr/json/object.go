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

func (o object) WithField(key string, value repr.Encodable) repr.Object {
	return object{fields: slices.With(o.fields, field{key: key, value: value})}
}

func Object() repr.Object {
	return object{fields: []field{}}
}

type field struct {
	key   string
	value repr.Encodable
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

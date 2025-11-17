package json

import (
	"bytes"
	"fmt"
	"maps"
	"slices"

	"github.com/andreychh/coopera-bot/pkg/repr"
)

type object struct {
	fields Fields
}

func (o object) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, key := range slices.Sorted(maps.Keys(o.fields)) {
		if i > 0 {
			buf.WriteByte(',')
		}
		keyData, err := String(key).Marshal()
		if err != nil {
			return nil, err
		}
		buf.Write(keyData)
		buf.WriteByte(':')
		valueData, err := o.fields[key].Marshal()
		if err != nil {
			return nil, err
		}
		buf.Write(valueData)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (o object) At(path repr.Path) (repr.Structure, error) {
	if path.Empty() {
		return o, nil
	}
	key, err := o.key(path)
	if err != nil {
		return nil, err
	}
	value, exists := o.fields[key]
	if !exists {
		return nil, fmt.Errorf("field %q not found", key)
	}
	return value.At(path.Tail())
}

func (o object) Update(path repr.Path, value repr.Structure) (repr.Structure, error) {
	if path.Empty() {
		return value, nil
	}
	key, err := o.key(path)
	if err != nil {
		return nil, err
	}
	existing, exists := o.fields[key]
	if !exists {
		return nil, fmt.Errorf("field %q not found", key)
	}
	updated, err := existing.Update(path.Tail(), value)
	if err != nil {
		return nil, err
	}
	newFields := make(Fields, len(o.fields))
	for k, v := range o.fields {
		newFields[k] = v
	}
	newFields[key] = updated
	return object{fields: newFields}, nil
}

func (o object) Extend(path repr.Path, other repr.Structure) (repr.Structure, error) {
	if path.Empty() {
		otherObject, ok := other.(object)
		if !ok {
			return nil, fmt.Errorf("can only extend object with object, got %T", other)
		}
		merged := make(Fields, len(o.fields)+len(otherObject.fields))
		for k, v := range o.fields {
			merged[k] = v
		}
		for k, v := range otherObject.fields {
			merged[k] = v
		}
		return object{fields: merged}, nil
	}
	key, err := o.key(path)
	if err != nil {
		return nil, err
	}
	existing, exists := o.fields[key]
	if !exists {
		return nil, fmt.Errorf("field %q not found", key)
	}
	extended, err := existing.Extend(path.Tail(), other)
	if err != nil {
		return nil, err
	}
	newFields := make(Fields, len(o.fields))
	for k, v := range o.fields {
		newFields[k] = v
	}
	newFields[key] = extended
	return object{fields: newFields}, nil
}

func (o object) key(path repr.Path) (string, error) {
	head, err := path.Head()
	if err != nil {
		return "", fmt.Errorf("getting path head: %w", err)
	}
	key, ok := head.Key()
	if !ok {
		return "", fmt.Errorf("object path must be key, got %v", head)
	}
	return key, nil
}

func Object(fields Fields) repr.Structure {
	copied := make(Fields, len(fields))
	for k, v := range fields {
		copied[k] = v
	}
	return object{fields: copied}
}

func EmptyObject() repr.Structure {
	return object{fields: make(Fields)}
}

type Fields map[string]repr.Structure

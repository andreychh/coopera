package keyvalue

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
)

type keyValueField struct {
	dataSource keyvalue.Store
	key        string
}

func (k keyValueField) Value(ctx context.Context) (string, error) {
	value, err := k.dataSource.Read(ctx, k.valueKey())
	if err != nil {
		return "", fmt.Errorf("(%T->%T) reading field(%s) value: %w", k, k.dataSource, k.key, err)
	}
	return value, nil
}

func (k keyValueField) ChangeValue(ctx context.Context, value string) error {
	err := k.dataSource.Write(ctx, k.valueKey(), value)
	if err != nil {
		return fmt.Errorf("(%T->%T) writing field(%s) value to %q: %w", k, k.dataSource, k.key, value, err)
	}
	return nil
}

func (k keyValueField) valueKey() string {
	return fmt.Sprintf("%s:value", k.key)
}

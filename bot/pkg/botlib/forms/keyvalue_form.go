package forms

import (
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
)

type KeyValueForm struct {
	dataSource keyvalue.Store
	key        string
}

func (k KeyValueForm) Field(name string) Field {
	return keyValueField{
		dataSource: k.dataSource,
		key:        k.fieldKey(name),
	}
}

func (k KeyValueForm) fieldKey(name string) string {
	return fmt.Sprintf("%s:field:%s", k.key, name)
}

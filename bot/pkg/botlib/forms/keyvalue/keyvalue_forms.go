package keyvalue

import (
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
)

type keyValueForms struct {
	dataSource keyvalue.Store
}

func (k keyValueForms) Form(id int64) forms.Form {
	return KeyValueForm{
		dataSource: k.dataSource,
		key:        k.formKey(id),
	}
}

func (k keyValueForms) formKey(id int64) string {
	return fmt.Sprintf("form:%d", id)
}

func KeyValueForms(dataSource keyvalue.Store) forms.Forms {
	return keyValueForms{dataSource: dataSource}
}

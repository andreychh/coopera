package forms

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
)

func FormWithID(forms Forms, id sources.Source[int64]) sources.Source[Form] {
	return sources.PureMap(id, forms.Form)
}

func CurrentForm(forms Forms) sources.Source[Form] {
	return FormWithID(forms, sources.Required(attributes.ChatID()))
}

func FieldWithName(form sources.Source[Form], name string) sources.Source[Field] {
	return sources.PureMap(form,
		func(f Form) Field {
			return f.Field(name)
		},
	)
}

func CurrentField(forms Forms, name string) sources.Source[Field] {
	return FieldWithName(CurrentForm(forms), name)
}

func Value(field sources.Source[Field]) sources.Source[string] {
	return sources.Map(field,
		func(ctx context.Context, f Field) (string, error) {
			return f.Value(ctx)
		},
	)
}

func CurrentValue(forms Forms, name string) sources.Source[string] {
	return Value(CurrentField(forms, name))
}

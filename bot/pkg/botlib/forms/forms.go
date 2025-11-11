package forms

import (
	"context"
)

type Forms interface {
	Form(id int64) Form
}

type Form interface {
	Field(name string) Field
}

type Field interface {
	Value(ctx context.Context) (string, error)
	ChangeValue(ctx context.Context, value string) error
}

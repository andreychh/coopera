package actions

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
)

func SaveTextToField(f forms.Forms, name string) core.Action {
	return sources.Apply(
		forms.CurrentField(f, name),
		sources.Required(attributes.Text()),
		func(ctx context.Context, field forms.Field, value string) error {
			return field.ChangeValue(ctx, value)
		},
	)
}

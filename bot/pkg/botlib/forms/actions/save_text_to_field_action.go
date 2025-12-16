package actions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

type safeValueToFieldAction struct {
	forms forms.Forms
	name  string
	value string
}

func (s safeValueToFieldAction) Perform(ctx context.Context, update telegram.Update) error {
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return fmt.Errorf("chat ID not found in update")
	}
	return s.forms.Form(chatID).Field(s.name).ChangeValue(ctx, s.value)
}

func SaveValueToField(f forms.Forms, name string, value string) core.Action {
	return safeValueToFieldAction{
		forms: f,
		name:  name,
		value: value,
	}
}

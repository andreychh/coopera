package actions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type saveTextToFieldAction struct {
	forms forms.Forms
	name  string
}

func (s saveTextToFieldAction) Perform(ctx context.Context, update telegram.Update) error {
	id, exists := attrs.ChatID(update).Value()
	if !exists {
		return fmt.Errorf("getting chat ID from update: chat ID not found")
	}
	text, exists := attrs.Text(update).Value()
	if !exists {
		return fmt.Errorf("getting text from update: text not found")
	}
	err := s.forms.Form(id).Field(s.name).ChangeValue(ctx, text)
	if err != nil {
		return fmt.Errorf("changing value of field %q in form for chat %d to %q: %w", s.name, id, text, err)
	}
	return nil
}

func SaveTextToField(forms forms.Forms, name string) core.Action {
	return saveTextToFieldAction{forms: forms, name: name}
}

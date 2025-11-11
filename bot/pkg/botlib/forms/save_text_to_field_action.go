package forms

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type saveTextToFieldAction struct {
	forms Forms
	name  string
}

func (s saveTextToFieldAction) Perform(ctx context.Context, update telegram.Update) error {
	id, available := updates.ChatID(update)
	if !available {
		return fmt.Errorf("(%T) getting chat ID: %w", s, updates.ErrNoChatID)
	}
	text, available := updates.Text(update)
	if !available {
		return fmt.Errorf("(%T) getting message text: %w", s, updates.ErrNoText)
	}
	err := s.forms.Form(id).Field(s.name).ChangeValue(ctx, text)
	if err != nil {
		return fmt.Errorf(
			"(%T->%T) saving text %q to field %q of form(%d): %w",
			s, s.forms, text, s.name, id, err,
		)
	}
	return nil
}

func SaveTextToField(forms Forms, name string) core.Action {
	return saveTextToFieldAction{forms: forms, name: name}
}

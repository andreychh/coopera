package dialogues

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type dialogueExistsCondition struct {
	dialogues Dialogues
}

func (d dialogueExistsCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	id, err := attributes.ChatID(update).Value()
	if err != nil {
		return false, fmt.Errorf("(%T) getting chat ID: %w", d, err)
	}
	exists, err := d.dialogues.Dialogue(id).Exists(ctx)
	if err != nil {
		return false, fmt.Errorf(
			"(%T->%T) checking dialogue existence for chat id #%d: %w",
			d, d.dialogues, id, err,
		)
	}
	return exists, nil
}

func DialogueExists(dialogues Dialogues) core.Condition {
	return dialogueExistsCondition{dialogues: dialogues}
}

func SafeDialogueExists(dialogues Dialogues) core.Condition {
	return composition.All(updates.HasChatID(), DialogueExists(dialogues))
}

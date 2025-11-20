package conditions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/dialogues"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type dialogueExistsCondition struct {
	dialogues dialogues.Dialogues
}

func (d dialogueExistsCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	id, exists := attrs.ChatID(update).Value()
	if !exists {
		return false, fmt.Errorf("getting chat ID from update: chat ID not found")
	}
	exists, err := d.dialogues.Dialogue(id).Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("checking existence of dialogue for chat %d: %w", id, err)
	}
	return exists, nil
}

func DialogueExists(dialogues dialogues.Dialogues) core.Condition {
	return dialogueExistsCondition{dialogues: dialogues}
}

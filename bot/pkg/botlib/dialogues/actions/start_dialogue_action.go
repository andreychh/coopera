package actions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/dialogues"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type startDialogueAction struct {
	dialogues dialogues.Dialogues
	topic     dialogues.Topic
}

func (s startDialogueAction) Perform(ctx context.Context, update telegram.Update) error {
	id, exists := attrs.ChatID(update).Value()
	if !exists {
		return fmt.Errorf("getting chat ID from update: chat ID not found")
	}
	err := s.dialogues.Dialogue(id).ChangeTopic(ctx, s.topic)
	if err != nil {
		return fmt.Errorf("changing topic of dialogue for chat %d to %q: %w", id, s.topic, err)
	}
	return nil
}

func StartDialogue(dialogues dialogues.Dialogues, topic dialogues.Topic) core.Action {
	return startDialogueAction{dialogues: dialogues, topic: topic}
}

func StartNeutralDialog(d dialogues.Dialogues) core.Action {
	return StartDialogue(d, dialogues.TopicNeutral)
}

package dialogues

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type startDialogueAction struct {
	dialogues Dialogues
	topic     Topic
}

func (s startDialogueAction) Perform(ctx context.Context, update telegram.Update) error {
	id, available := updates.ChatID(update)
	if !available {
		return fmt.Errorf("(%T) getting chat ID: %w", s, updates.ErrNoChatID)
	}
	_, err := s.dialogues.StartDialogue(ctx, id, s.topic)
	if err != nil {
		return fmt.Errorf(
			"(%T->%T) starting dialogue for chat id #%d with topic %s: %w",
			s, s.dialogues, id, s.topic, err,
		)
	}
	return nil
}

func StartDialogue(dialogues Dialogues, topic Topic) core.Action {
	return startDialogueAction{dialogues: dialogues, topic: topic}
}

func StartNeutralDialog(dialogues Dialogues) core.Action {
	return StartDialogue(dialogues, TopicNeutral)
}

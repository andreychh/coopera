package dialogues

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type changeTopicAction struct {
	dialogues Dialogues
	topic     Topic
}

func (c changeTopicAction) Perform(ctx context.Context, update telegram.Update) error {
	id, available := updates.ChatID(update)
	if !available {
		return fmt.Errorf("(%T) getting chat ID: %w", c, updates.ErrNoChatID)
	}
	err := c.dialogues.Dialogue(id).ChangeTopic(ctx, c.topic)
	if err != nil {
		return fmt.Errorf("(%T->%T) changing topic for chat id #%d to %s: %w", c, c.dialogues, id, c.topic, err)
	}
	return nil
}

func ChangeTopic(dialogues Dialogues, topic Topic) core.Action {
	return changeTopicAction{dialogues: dialogues, topic: topic}
}

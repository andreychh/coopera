package actions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/dialogues"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type changeTopicAction struct {
	dialogues dialogues.Dialogues
	topic     dialogues.Topic
}

func (c changeTopicAction) Perform(ctx context.Context, update telegram.Update) error {
	id, exists := attrs.ChatID(update).Value()
	if !exists {
		return fmt.Errorf("getting chat ID from update: chat ID not found")
	}
	err := c.dialogues.Dialogue(id).ChangeTopic(ctx, c.topic)
	if err != nil {
		return fmt.Errorf("changing topic of dialogue for chat %d to %q: %w", id, c.topic, err)
	}
	return nil
}

func ChangeTopic(dialogues dialogues.Dialogues, topic dialogues.Topic) core.Action {
	return changeTopicAction{dialogues: dialogues, topic: topic}
}

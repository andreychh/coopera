package conditions

import (
	"context"

	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/dialogues"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type topicIsCondition struct {
	dialogues dialogues.Dialogues
	topic     dialogues.Topic
}

func (t topicIsCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	id, exists := attrs.ChatID(update).Value()
	if !exists {
		return false, fmt.Errorf("getting chat ID from update: chat ID not found")
	}
	topic, err := t.dialogues.Dialogue(id).Topic(ctx)
	if err != nil {
		return false, fmt.Errorf("getting topic of dialogue for chat %d: %w", id, err)
	}
	return topic == t.topic, nil
}

func TopicIs(dialogues dialogues.Dialogues, topic dialogues.Topic) core.Condition {
	return topicIsCondition{dialogues: dialogues, topic: topic}
}

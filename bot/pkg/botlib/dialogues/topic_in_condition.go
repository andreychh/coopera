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

type topicInCondition struct {
	dialogues Dialogues
	topics    []Topic
}

func (t topicInCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	id, err := attributes.ChatID(update).Value()
	if err != nil {
		return false, fmt.Errorf("(%T) getting chat ID: %w", t, err)
	}
	topic, err := t.dialogues.Dialogue(id).Topic(ctx)
	if err != nil {
		return false, fmt.Errorf("(%T->%T) getting topic for chat id #%d: %w", t, t.dialogues, id, err)
	}
	for _, tpc := range t.topics {
		if tpc == topic {
			return true, nil
		}
	}
	return false, nil
}

func TopicIn(dialogues Dialogues, topics ...Topic) core.Condition {
	return topicInCondition{dialogues: dialogues, topics: topics}
}

func SafeTopicIn(dialogues Dialogues, topics ...Topic) core.Condition {
	return composition.All(updates.HasChatID(), SafeDialogueExists(dialogues), TopicIn(dialogues, topics...))
}

func TopicIs(dialogues Dialogues, topic Topic) core.Condition {
	return TopicIn(dialogues, topic)
}

func SafeTopicIs(dialogues Dialogues, topic Topic) core.Condition {
	return SafeTopicIn(dialogues, topic)
}

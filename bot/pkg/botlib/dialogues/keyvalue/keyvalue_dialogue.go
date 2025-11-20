package keyvalue

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/dialogues"
	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
)

type keyValueDialogue struct {
	dataSource keyvalue.Store
	key        string
}

func (k keyValueDialogue) ChangeTopic(ctx context.Context, topic dialogues.Topic) error {
	err := k.dataSource.Write(ctx, k.topicKey(), string(topic))
	if err != nil {
		return fmt.Errorf(
			"(%T->%T) writing dialogue(%s) topic to %q: %w",
			k, k.dataSource, k.key, topic, err,
		)
	}
	return nil
}

func (k keyValueDialogue) Topic(ctx context.Context) (dialogues.Topic, error) {
	topic, err := k.dataSource.Read(ctx, k.topicKey())
	if err != nil {
		return dialogues.TopicNeutral, fmt.Errorf(
			"(%T->%T) reading dialogue(%s) topic: %w",
			k, k.dataSource, k.key, err,
		)
	}
	return dialogues.Topic(topic), nil
}

func (k keyValueDialogue) Exists(ctx context.Context) (bool, error) {
	exists, err := k.dataSource.Exists(ctx, k.topicKey())
	if err != nil {
		return false, fmt.Errorf(
			"(%T->%T) checking existence of dialogue(%s) topic: %w",
			k, k.dataSource, k.key, err,
		)
	}
	return exists, nil
}

func (k keyValueDialogue) topicKey() string {
	return fmt.Sprintf("%s:topic", k.key)
}

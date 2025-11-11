package dialogues

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
)

type keyValueDialogue struct {
	id         int64
	dataSource keyvalue.Store
}

func (k keyValueDialogue) ChangeTopic(ctx context.Context, topic Topic) error {
	err := k.dataSource.Write(ctx, k.key(), string(topic))
	if err != nil {
		return fmt.Errorf("(%T->%T) writing topic for chat id #%d: %w", k, k.dataSource, k.id, err)
	}
	return nil
}

func (k keyValueDialogue) Topic(ctx context.Context) (Topic, error) {
	topic, err := k.dataSource.Read(ctx, k.key())
	if err != nil {
		return TopicNeutral, fmt.Errorf("(%T->%T) reading topic for chat id #%d: %w", k, k.dataSource, k.id, err)
	}
	return Topic(topic), nil
}

func (k keyValueDialogue) Exists(ctx context.Context) (bool, error) {
	exists, err := k.dataSource.Exists(ctx, k.key())
	if err != nil {
		return false, fmt.Errorf("(%T->%T) checking existence for chat id #%d: %w", k, k.dataSource, k.id, err)
	}
	return exists, nil
}

func (k keyValueDialogue) key() string {
	return fmt.Sprintf("chat:%d:dialogue", k.id)
}

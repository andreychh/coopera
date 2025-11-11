package dialogues

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
)

type keyValueDialogue struct {
	key        string
	dataSource keyvalue.Store
}

func (k keyValueDialogue) ChangeTopic(ctx context.Context, topic Topic) error {
	err := k.dataSource.Write(ctx, k.key, string(topic))
	if err != nil {
		return fmt.Errorf("(%T->%T) changing topic for dialogue #%s: %w", k, k.dataSource, k.key, err)
	}
	return nil
}

func (k keyValueDialogue) Topic(ctx context.Context) (Topic, error) {
	topic, err := k.dataSource.Read(ctx, k.key)
	if err != nil {
		return TopicNeutral, fmt.Errorf("(%T->%T) getting topic for dialogue #%s: %w", k, k.dataSource, k.key, err)
	}
	return Topic(topic), nil
}

func (k keyValueDialogue) Exists(ctx context.Context) (bool, error) {
	exists, err := k.dataSource.Exists(ctx, k.key)
	if err != nil {
		return false, fmt.Errorf("(%T->%T) checking existence for dialogue #%s: %w", k, k.dataSource, k.key, err)
	}
	return exists, nil
}

package dialogues

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
)

type keyValueDialogues struct {
	dataSource keyvalue.Store
}

func (k keyValueDialogues) StartDialogue(ctx context.Context, id int64, topic Topic) (Dialogue, error) {
	err := k.dataSource.Write(ctx, k.key(id), string(topic))
	if err != nil {
		return nil, fmt.Errorf(
			"(%T->%T) starting dialogue for chat id #%d with topic %s: %w",
			k, k.dataSource, id, topic, err,
		)
	}
	return k.Dialogue(id), nil
}

func (k keyValueDialogues) Dialogue(id int64) Dialogue {
	return keyValueDialogue{
		key:        k.key(id),
		dataSource: k.dataSource,
	}
}

func (k keyValueDialogues) key(id int64) string {
	return fmt.Sprintf("chat:%d:dialogue", id)
}

func KeyValueDialogues(dataSource keyvalue.Store) Dialogues {
	return keyValueDialogues{dataSource: dataSource}
}

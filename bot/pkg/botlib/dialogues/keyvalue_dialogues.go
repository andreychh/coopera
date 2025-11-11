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
	dialogue := k.Dialogue(id)
	err := dialogue.ChangeTopic(ctx, topic)
	if err != nil {
		return nil, fmt.Errorf("(%T->%T) starting dialogue for chat id #%d with topic %s: %w", k, k.dataSource, id, topic, err)
	}
	return dialogue, nil
}

func (k keyValueDialogues) Dialogue(id int64) Dialogue {
	return keyValueDialogue{
		id:         id,
		dataSource: k.dataSource,
	}
}

func KeyValueDialogues(dataSource keyvalue.Store) Dialogues {
	return keyValueDialogues{dataSource: dataSource}
}

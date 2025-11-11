package dialogues

import (
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
)

type keyValueDialogues struct {
	dataSource keyvalue.Store
}

func (k keyValueDialogues) Dialogue(id int64) Dialogue {
	return keyValueDialogue{
		dataSource: k.dataSource,
		key:        k.dialogueKey(id),
	}
}

func (k keyValueDialogues) dialogueKey(id int64) string {
	return fmt.Sprintf("dialogue:%d", id)
}

func KeyValueDialogues(dataSource keyvalue.Store) Dialogues {
	return keyValueDialogues{dataSource: dataSource}
}

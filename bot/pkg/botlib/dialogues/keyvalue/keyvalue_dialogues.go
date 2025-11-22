package keyvalue

import (
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/dialogues"
	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
)

type keyValueDialogues struct {
	dataSource keyvalue.Store
}

func (k keyValueDialogues) Dialogue(id int64) dialogues.Dialogue {
	return keyValueDialogue{
		dataSource: k.dataSource,
		key:        k.dialogueKey(id),
	}
}

func (k keyValueDialogues) dialogueKey(id int64) string {
	return fmt.Sprintf("dialogue:%d", id)
}

func KeyValueDialogues(dataSource keyvalue.Store) dialogues.Dialogues {
	return keyValueDialogues{dataSource: dataSource}
}

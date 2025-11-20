package conditions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type commandIsCondition struct {
	target string
}

func (c commandIsCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	command, exists := attrs.Command(update).Value()
	if !exists {
		return false, fmt.Errorf("getting command from update: command not found")
	}
	return command == c.target, nil
}

func CommandIs(target string) core.Condition {
	return commandIsCondition{target: target}
}

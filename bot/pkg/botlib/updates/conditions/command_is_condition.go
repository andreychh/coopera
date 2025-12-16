package conditions

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type commandIsCondition struct {
	target string
}

func (c commandIsCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	command, exists := attributes.Command().Value(update)
	if !exists {
		return false, nil
	}
	return command == c.target, nil
}

func CommandIs(target string) core.Condition {
	return commandIsCondition{target: target}
}

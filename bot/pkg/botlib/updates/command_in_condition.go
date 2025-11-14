package updates

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type commandInCondition struct {
	commands []string
}

func (a commandInCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	command, err := attributes.Command(update).Value()
	if err != nil {
		return false, fmt.Errorf("(%T) getting command: %w", a, err)
	}
	for _, cmd := range a.commands {
		if cmd == command {
			return true, nil
		}
	}
	return false, nil
}

func CommandIn(commands ...string) core.Condition {
	return commandInCondition{commands: commands}
}

func SafeCommandIn(commands ...string) core.Condition {
	return composition.All(HasCommand(), CommandIn(commands...))
}

func CommandIs(command string) core.Condition {
	return CommandIn(command)
}

func SafeCommandIs(command string) core.Condition {
	return SafeCommandIn(command)
}

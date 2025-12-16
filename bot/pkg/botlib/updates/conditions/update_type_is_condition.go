package conditions

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type updateTypeIsCondition struct {
	target updates.UpdateType
}

func (u updateTypeIsCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	updateType, exists := attributes.UpdateType().Value(update)
	if !exists {
		return false, nil
	}
	return updateType == u.target, nil
}

func UpdateTypeIs(target updates.UpdateType) core.Condition {
	return updateTypeIsCondition{target: target}
}

package conditions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type updateTypeIsCondition struct {
	updateType updates.UpdateType
}

func (u updateTypeIsCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	updateType, exists := attrs.UpdateType(update).Value()
	if !exists {
		return false, fmt.Errorf("getting update type: update type not found")
	}
	return updateType == u.updateType, nil
}

func UpdateTypeIs(updateType updates.UpdateType) core.Condition {
	return updateTypeIsCondition{updateType: updateType}
}

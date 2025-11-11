package updates

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type hasChatIDCondition struct{}

func (c hasChatIDCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	_, available := ChatID(update)
	return available, nil
}

func HasChatID() core.Condition {
	return hasChatIDCondition{}
}

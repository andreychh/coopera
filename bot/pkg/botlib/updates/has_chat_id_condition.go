package updates

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type hasChatIDCondition struct{}

func (c hasChatIDCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	return attributes.ChatID(update).Exists(), nil
}

func HasChatID() core.Condition {
	return hasChatIDCondition{}
}

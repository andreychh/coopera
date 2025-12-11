package base

import (
	"context"
	"errors"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type editOrSendContentAction struct {
	bot     tg.Bot
	content sources.Source[content.Content]
}

func (e editOrSendContentAction) Perform(ctx context.Context, update telegram.Update) error {
	updateType, exists := attributes.UpdateType().Value(update)
	if !exists {
		return fmt.Errorf("cannot determine update type")
	}
	if updateType == updates.UpdateTypeCallbackQuery {
		err := EditMessage(e.bot, e.content).Perform(ctx, update)
		if err == nil {
			return nil
		}
		if !errors.Is(err, tg.ErrMessageCannotBeEdited) && !errors.Is(err, tg.ErrMessageNotFound) {
			return fmt.Errorf("editing message failed: %w", err)
		}
	}
	return SendContent(e.bot, e.content).Perform(ctx, update)
}

func EditOrSendContent(bot tg.Bot, content sources.Source[content.Content]) core.Action {
	return editOrSendContentAction{
		bot:     bot,
		content: content,
	}
}

package actions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type markTaskAsCompletedAction struct {
	community domain.Community
}

func (m markTaskAsCompletedAction) Perform(ctx context.Context, update telegram.Update) error {
	callback, found := attributes.CallbackData().Value(update)
	if !found {
		return fmt.Errorf("callback data not found in update")
	}
	taskID, err := protocol.ParseTaskID(callback)
	if err != nil {
		return fmt.Errorf("parsing task ID from callback data %q: %w", callback, err)
	}
	task, err := m.community.Task(ctx, taskID)
	if err != nil {
		return fmt.Errorf("getting task %d: %w", taskID, err)
	}
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return fmt.Errorf("chat ID not found in update")
	}
	user, err := m.community.UserWithTelegramID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("getting user with telegram ID %d: %w", chatID, err)
	}
	err = task.MarkAsCompleted(ctx, user)
	if err != nil {
		return fmt.Errorf("marking task %d as completed by user %d: %w", taskID, user.ID(), err)
	}
	return nil
}

func MarkTaskAsCompleted(community domain.Community) core.Action {
	return markTaskAsCompletedAction{
		community: community,
	}
}

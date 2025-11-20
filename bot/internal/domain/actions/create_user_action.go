package actions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type createUserAction struct {
	community domain.Community
}

func (c createUserAction) Perform(ctx context.Context, update telegram.Update) error {
	id, exists := attrs.ChatID(update).Value()
	if !exists {
		return fmt.Errorf("getting chat ID from update: chat ID not found")
	}
	_, err := c.community.CreateUser(ctx, id)
	if err != nil {
		return fmt.Errorf("creating user for chat %d: %w", id, err)
	}
	return nil
}

func CreateUser(community domain.Community) core.Action {
	return createUserAction{community: community}
}

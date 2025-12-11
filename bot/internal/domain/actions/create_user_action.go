package actions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type createUserAction struct {
	community domain.Community
}

func (c createUserAction) Perform(ctx context.Context, update telegram.Update) error {
	id, found := attributes.ChatID().Value(update)
	if !found {
		return fmt.Errorf("chat ID not found in update")
	}
	username, found := attributes.Username().Value(update)
	if !found {
		return fmt.Errorf("username not found in update")
	}
	_, err := c.community.CreateUser(ctx, id, username)
	if err != nil {
		return fmt.Errorf("creating user %d (%s): %w", id, username, err)
	}
	return nil
}

func CreateUser(community domain.Community) core.Action {
	return createUserAction{
		community: community,
	}
}

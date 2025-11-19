package domain

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type createUserAction struct {
	community Community
}

func (c createUserAction) Perform(ctx context.Context, update telegram.Update) error {
	id, err := attributes.ChatID(update).Value()
	if err != nil {
		return fmt.Errorf("(%T) getting chat ID: %w", c, err)
	}
	_, err = c.community.CreateUser(ctx, id)
	if err != nil {
		return fmt.Errorf("(%T) creating user with telegram ID #%d: %w", c, id, err)
	}
	return nil
}

func CreateUser(community Community) core.Action {
	return createUserAction{community: community}
}

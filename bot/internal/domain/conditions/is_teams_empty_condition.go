package conditions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type isTeamsEmptyCondition struct {
	community domain.Community
}

func (i isTeamsEmptyCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	id, found := attributes.ChatID().Value(update)
	if !found {
		return false, fmt.Errorf("chat ID not found in update")
	}
	user, exists, err := i.community.UserWithTelegramID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("getting user with telegram ID %d: %w", id, err)
	}
	if !exists {
		return false, fmt.Errorf("user with telegram ID %d does not exist", id)
	}
	teams, err := user.Teams(ctx)
	if err != nil {
		return false, fmt.Errorf("getting created teams for user %d: %w", id, err)
	}
	isEmpty, err := teams.Empty(ctx)
	if err != nil {
		return false, fmt.Errorf("checking if teams are empty for user %d: %w", id, err)
	}
	return isEmpty, nil
}

func IsTeamsEmpty(community domain.Community) core.Condition {
	return isTeamsEmptyCondition{
		community: community,
	}
}

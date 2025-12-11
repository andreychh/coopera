package conditions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type teamExistsCondition struct {
	community domain.Community
}

func (t teamExistsCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	id, found := attributes.ChatID().Value(update)
	if !found {
		return false, fmt.Errorf("chat ID not found in update")
	}
	name, found := attributes.Text().Value(update)
	if !found {
		return false, fmt.Errorf("text not found in update")
	}
	user, err := t.community.UserWithTelegramID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("getting user with telegram ID %d: %w", id, err)
	}
	teams, err := user.CreatedTeams(ctx)
	if err != nil {
		return false, fmt.Errorf("getting created teams for user %d: %w", id, err)
	}
	exists, err := teams.ContainsTeam(ctx, name)
	if err != nil {
		return false, fmt.Errorf("checking if team %q exists for user %d: %w", name, id, err)
	}
	return exists, nil
}

func TeamExists(community domain.Community) core.Condition {
	return teamExistsCondition{
		community: community,
	}
}

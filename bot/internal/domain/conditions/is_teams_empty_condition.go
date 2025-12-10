package conditions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type isTeamsEmptyCondition struct {
	community domain.Community
	id        sources.Source[int64]
}

func (i isTeamsEmptyCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	id, err := i.id.Value(ctx, update)
	if err != nil {
		return false, fmt.Errorf("getting user ID: %w", err)
	}
	isEmpty, err := i.community.UserWithTelegramID(id).CreatedTeams().Empty(ctx)
	if err != nil {
		return false, fmt.Errorf("checking if teams are empty: %w", err)
	}
	return isEmpty, nil
}

func IsTeamsEmpty(community domain.Community) core.Condition {
	return isTeamsEmptyCondition{
		community: community,
		id:        sources.Required(attributes.ChatID()),
	}
}

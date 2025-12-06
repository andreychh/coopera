package conditions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type teamExistsCondition struct {
	community domain.Community
}

func (t teamExistsCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	id, exists := attrs.ChatID(update).Value()
	if !exists {
		return false, fmt.Errorf("getting chat ID from update: chat ID not found")
	}
	text, exists := attrs.Text(update).Value()
	if !exists {
		return false, fmt.Errorf("getting text from update: text not found")
	}
	teams, err := t.community.UserWithTelegramID(id).CreatedTeams(ctx)
	if err != nil {
		return false, err
	}
	for _, team := range teams {
		details, err := team.Details(ctx)
		if err != nil {
			return false, err
		}
		if details.Name() == text {
			return true, nil
		}
	}
	return false, nil
}

func TeamExists(community domain.Community) core.Condition {
	return teamExistsCondition{
		community: community,
	}
}

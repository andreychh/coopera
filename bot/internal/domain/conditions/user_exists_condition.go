package conditions

import (
	"context"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type userExistsCondition struct {
	community domain.Community
}

func (u userExistsCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	username, found := attributes.Text().Value(update)
	if !found {
		return false, nil
	}
	return u.community.UserWithUsernameExists(ctx, username)
}

func UserExists(community domain.Community) core.Condition {
	return userExistsCondition{
		community: community,
	}
}

package conditions

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type userAlreadyInTeamCondition struct {
	community domain.Community
	forms     forms.Forms
}

func (u userAlreadyInTeamCondition) Holds(ctx context.Context, update telegram.Update) (bool, error) {
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return false, fmt.Errorf("chat ID not found in update")
	}
	teamID, err := u.forms.Form(chatID).Field("team_id").Value(ctx)
	if err != nil {
		return false, fmt.Errorf("getting team_id field for user %d: %w", chatID, err)
	}
	intTeamID, err := strconv.ParseInt(teamID, 10, 64)
	if err != nil {
		return false, fmt.Errorf("parsing team_id %q to int64: %w", teamID, err)
	}
	username, found := attributes.Text().Value(update)
	if !found {
		return false, nil
	}
	team, err := u.community.Team(ctx, intTeamID)
	if err != nil {
		return false, fmt.Errorf("getting team with ID %d: %w", intTeamID, err)
	}
	user, err := u.community.UserWithUsername(ctx, username)
	if err != nil {
		return false, fmt.Errorf("getting user with username %q: %w", username, err)
	}
	return team.ContainsUser(ctx, user)
}

func UserAlreadyInTeam(community domain.Community, forms forms.Forms) core.Condition {
	return userAlreadyInTeamCondition{
		community: community,
		forms:     forms,
	}
}

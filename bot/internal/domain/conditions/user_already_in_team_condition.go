package conditions

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
	teamIDStr, err := u.forms.Form(chatID).Field("team_id").Value(ctx)
	if err != nil {
		return false, fmt.Errorf("getting team_id field for user %d: %w", chatID, err)
	}
	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		return false, fmt.Errorf("parsing team_id %d to int64: %w", teamID, err)
	}
	username, found := attributes.Text().Value(update)
	if !found {
		return false, nil
	}
	team, exists, err := u.community.Team(ctx, teamID)
	if err != nil {
		return false, fmt.Errorf("getting team with ID %d: %w", teamID, err)
	}
	if !exists {
		return false, fmt.Errorf("team with ID %d does not exist", teamID)
	}
	username = strings.TrimPrefix(username, "@")
	user, exists, err := u.community.UserWithUsername(ctx, username)
	if err != nil {
		return false, fmt.Errorf("getting user with username %q: %w", username, err)
	}
	if !exists {
		return false, nil
	}
	members, err := team.Members(ctx)
	if err != nil {
		return false, fmt.Errorf("getting members of team %d: %w", teamID, err)
	}
	_, exists, err = members.MemberWithUsername(ctx, user.Username())
	if err != nil {
		return false, fmt.Errorf("checking if user %q is member of team %d: %w", username, teamID, err)
	}
	return exists, nil
}

func UserInTeam(community domain.Community, forms forms.Forms) core.Condition {
	return userAlreadyInTeamCondition{
		community: community,
		forms:     forms,
	}
}

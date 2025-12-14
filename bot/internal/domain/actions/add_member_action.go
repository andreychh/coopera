package actions

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

type addMemberAction struct {
	community domain.Community
	forms     forms.Forms
}

func (a addMemberAction) Perform(ctx context.Context, update telegram.Update) error {
	chatID, found := attributes.ChatID().Value(update)
	if !found {
		return fmt.Errorf("chat ID not found in update")
	}
	teamID, err := a.forms.Form(chatID).Field("team_id").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting team_id field for user %d: %w", chatID, err)
	}
	intTeamID, err := strconv.ParseInt(teamID, 10, 64)
	if err != nil {
		return fmt.Errorf("parsing team_id %q to int64: %w", teamID, err)
	}
	username, err := a.forms.Form(chatID).Field("member_username").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting member_username field for user %d: %w", chatID, err)
	}
	user, exists, err := a.community.UserWithUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("getting user with username %q: %w", username, err)
	}
	if !exists {
		return fmt.Errorf("user with username %q does not exist", username)
	}
	team, exists, err := a.community.Team(ctx, intTeamID)
	if err != nil {
		return fmt.Errorf("getting team with ID %d: %w", intTeamID, err)
	}
	if !exists {
		return fmt.Errorf("team with ID %d does not exist", intTeamID)
	}
	_, err = team.AddMember(ctx, user.ID())
	if err != nil {
		return fmt.Errorf("adding user %q to team %d: %w", username, intTeamID, err)
	}
	return nil
}

func AddMember(community domain.Community, forms forms.Forms) core.Action {
	return addMemberAction{
		community: community,
		forms:     forms,
	}
}

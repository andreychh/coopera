package actions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type createTeamAction struct {
	community domain.Community
	forms     forms.Forms
}

func (c createTeamAction) Perform(ctx context.Context, update telegram.Update) error {
	id, found := attributes.ChatID().Value(update)
	if !found {
		return fmt.Errorf("chat ID not found in update")
	}
	name, err := c.forms.Form(id).Field("team_name").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting team_name field for user %d: %w", id, err)
	}
	user, err := c.community.UserWithTelegramID(ctx, id)
	if err != nil {
		return fmt.Errorf("getting user with telegram ID %d: %w", id, err)
	}
	_, err = user.CreateTeam(ctx, name)
	if err != nil {
		return fmt.Errorf("creating team %q for user %d: %w", name, id, err)
	}
	return nil
}

func CreateTeam(community domain.Community, forms forms.Forms) core.Action {
	return createTeamAction{
		community: community,
		forms:     forms,
	}
}

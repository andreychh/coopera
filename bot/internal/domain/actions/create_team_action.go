package actions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type createTeamAction struct {
	forms     forms.Forms
	community domain.Community
}

func (c createTeamAction) Perform(ctx context.Context, update telegram.Update) error {
	id, exists := attrs.ChatID(update).Value()
	if !exists {
		return fmt.Errorf("getting chat ID: chat ID not found")
	}
	name, err := c.forms.Form(id).Field("team_name").Value(ctx)
	if err != nil {
		return fmt.Errorf("getting team name from form for chat %d: %w", id, err)
	}
	_, err = c.community.UserWithTelegramID(id).CreateTeam(ctx, name)
	if err != nil {
		return fmt.Errorf("creating team for user with chat ID %d: %w", id, err)
	}
	return nil
}

func CreateTeam(forms forms.Forms, community domain.Community) core.Action {
	return createTeamAction{
		forms:     forms,
		community: community,
	}
}

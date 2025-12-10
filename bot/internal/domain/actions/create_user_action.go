package actions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
)

func CreateUser(community domain.Community) core.Action {
	return sources.Apply(
		sources.Required(attributes.ChatID()),
		sources.Required(attributes.Username()),
		func(ctx context.Context, id int64, username string) error {
			_, err := community.CreateUser(ctx, id, username)
			if err != nil {
				return fmt.Errorf("creating user %d (%s): %w", id, username, err)
			}
			return nil
		},
	)
}

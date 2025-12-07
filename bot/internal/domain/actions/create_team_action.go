package actions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/internal/domain"
	sources2 "github.com/andreychh/coopera-bot/internal/domain/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/forms"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
)

func CreateTeam(f forms.Forms, community domain.Community) core.Action {
	return sources.Apply(
		sources2.CurrentUser(community),
		forms.CurrentValue(f, "team_name"),
		func(ctx context.Context, user domain.User, name string) error {
			_, err := user.CreateTeam(ctx, name)
			if err != nil {
				return fmt.Errorf("creating team %q: %w", name, err)
			}
			return nil
		},
	)
}

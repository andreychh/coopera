package base

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
)

func EditMessage(bot tg.Bot, cnt sources.Source[content.Content]) core.Action {
	return sources.Apply(
		tg.CurrentMessage(bot),
		cnt,
		func(ctx context.Context, message tg.Message, content content.Content) error {
			return message.Edit(ctx, content)
		},
	)
}

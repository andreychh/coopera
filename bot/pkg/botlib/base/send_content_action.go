package base

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/tg"
)

func SendContent(bot tg.Bot, cnt sources.Source[content.Content]) core.Action {
	return sources.Apply(
		tg.CurrentChat(bot),
		cnt,
		func(ctx context.Context, chat tg.Chat, content content.Content) error {
			return chat.Send(ctx, content)
		},
	)
}

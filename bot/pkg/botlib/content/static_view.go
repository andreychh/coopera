package content

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type staticView struct {
	content Content
}

func (s staticView) Render(_ context.Context, _ telegram.Update) (Content, error) {
	return s.content, nil
}

func StaticView(content Content) View {
	return staticView{
		content: content,
	}
}

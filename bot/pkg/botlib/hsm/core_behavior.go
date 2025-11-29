package hsm

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type coreBehavior struct {
	prompt  core.Action
	react   Handler
	cleanup core.Action
}

func (b coreBehavior) Prompt(ctx context.Context, u telegram.Update) error {
	return b.prompt.Perform(ctx, u)
}

func (b coreBehavior) React(ctx context.Context, u telegram.Update) (Decision, error) {
	return b.react.Handle(ctx, u)
}

func (b coreBehavior) Cleanup(ctx context.Context, u telegram.Update) error {
	return b.cleanup.Perform(ctx, u)
}

func CoreBehavior(enter core.Action, handle Handler, exit core.Action) Behavior {
	return coreBehavior{
		prompt:  enter,
		react:   DefaultLogging(handle),
		cleanup: exit,
	}
}

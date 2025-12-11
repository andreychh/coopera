package hsm

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Transition struct {
	targetID core.ID
	exit     path
	enter    path
}

func (t Transition) TargetID() core.ID {
	return t.targetID
}

func (t Transition) PerformCleanup(ctx context.Context, u telegram.Update) error {
	for state := range t.exit.All() {
		err := state.Exit(ctx, u)
		if err != nil {
			slog.Warn("exit hook failed", "state", state.ID(), "error", err)
		}
	}
	return nil
}

func (t Transition) PerformPrompt(ctx context.Context, u telegram.Update) error {
	for state := range t.enter.All() {
		err := state.Enter(ctx, u)
		if err != nil {
			return fmt.Errorf("enter hook failed at state %q: %w", state.ID(), err)
		}
	}
	return nil
}

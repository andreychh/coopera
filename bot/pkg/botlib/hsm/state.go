package hsm

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type State interface {
	ID() core.ID
	Parent() (State, bool)
	Initial() (State, bool)
	Enter(ctx context.Context, u telegram.Update) error
	Handle(ctx context.Context, u telegram.Update) (Decision, error)
	Exit(ctx context.Context, u telegram.Update) error
}

type generalState struct {
	id       core.ID
	behavior Behavior
	parent   State
	initial  State
}

func (s *generalState) ID() core.ID { return s.id }

func (s *generalState) Parent() (State, bool) {
	if s.parent == EdgeState() {
		return EdgeState(), false
	}
	return s.parent, true
}

func (s *generalState) Initial() (State, bool) {
	if s.initial == EdgeState() {
		return EdgeState(), false
	}
	return s.initial, true
}

func (s *generalState) Handle(ctx context.Context, u telegram.Update) (Decision, error) {
	return s.behavior.React(ctx, u)
}
func (s *generalState) Enter(ctx context.Context, u telegram.Update) error {
	return s.behavior.Prompt(ctx, u)
}
func (s *generalState) Exit(ctx context.Context, u telegram.Update) error {
	return s.behavior.Cleanup(ctx, u)
}

type edgeState struct{}

func (n edgeState) ID() core.ID {
	return core.Stay
}

func (n edgeState) Parent() (State, bool) {
	return nil, false
}

func (n edgeState) Initial() (State, bool) {
	return nil, false
}

func (n edgeState) Enter(ctx context.Context, u telegram.Update) error {
	return nil
}

func (n edgeState) Handle(ctx context.Context, u telegram.Update) (Decision, error) {
	return Stay(), nil
}

func (n edgeState) Exit(ctx context.Context, u telegram.Update) error {
	return nil
}

func EdgeState() State {
	return edgeState{}
}

package sources

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type appliedAction[A, B any] struct {
	a           Source[A]
	b           Source[B]
	application Application[A, B]
}

func (aa appliedAction[A, B]) Perform(ctx context.Context, update telegram.Update) error {
	a, err := aa.a.Value(ctx, update)
	if err != nil {
		return fmt.Errorf("getting A value: %w", err)
	}
	b, err := aa.b.Value(ctx, update)
	if err != nil {
		return fmt.Errorf("getting B value: %w", err)
	}
	return aa.application.Apply(ctx, a, b)
}

func Applied[A, B any](a Source[A], b Source[B], application Application[A, B]) core.Action {
	return appliedAction[A, B]{
		a:           a,
		b:           b,
		application: application,
	}
}

func Apply[A, B any](a Source[A], b Source[B], application func(context.Context, A, B) error) core.Action {
	return Applied(a, b, ApplicationFunc(
		func(ctx context.Context, a A, b B) error {
			return application(ctx, a, b)
		},
	))
}

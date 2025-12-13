package conditions

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type textIsCondition struct {
	expected string
}

func (t textIsCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	text, exists := attributes.Text().Value(update)
	if !exists {
		return false, nil
	}
	return text == t.expected, nil
}

func TextIs(expected string) core.Condition {
	return textIsCondition{
		expected: expected,
	}
}

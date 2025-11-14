package updates

import (
	"context"
	"fmt"
	"regexp"

	"github.com/andreychh/coopera-bot/pkg/botlib/composition"
	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type textMatchesRegexpCondition struct {
	pattern string
}

func (r textMatchesRegexpCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	text, err := attributes.Text(update).Value()
	if err != nil {
		return false, fmt.Errorf("(%T) getting message text: %w", r, err)
	}
	matched, err := regexp.MatchString(r.pattern, text)
	if err != nil {
		return false, fmt.Errorf("(%T) matching text against pattern %q: %w", r, r.pattern, err)
	}
	return matched, nil
}

func TextMatchesRegexp(pattern string) core.Condition {
	return textMatchesRegexpCondition{pattern: pattern}
}

func SafeTextMatchesRegexp(pattern string) core.Condition {
	return composition.All(HasText(), TextMatchesRegexp(pattern))
}

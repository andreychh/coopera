package conditions

import (
	"context"
	"regexp"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type textMatchesRegexpCondition struct {
	re *regexp.Regexp
}

func (r textMatchesRegexpCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	text, exists := attributes.Text().Value(update)
	if !exists {
		return false, nil
	}
	return r.re.MatchString(text), nil
}

func TextMatchesRegexp(pattern string) core.Condition {
	return textMatchesRegexpCondition{
		re: regexp.MustCompile(pattern),
	}
}

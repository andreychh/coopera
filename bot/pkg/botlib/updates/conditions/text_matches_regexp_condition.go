package conditions

import (
	"context"
	"fmt"
	"regexp"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type textMatchesRegexpCondition struct {
	re *regexp.Regexp
}

func (r textMatchesRegexpCondition) Holds(_ context.Context, update telegram.Update) (bool, error) {
	text, exists := attrs.Text(update).Value()
	if !exists {
		return false, fmt.Errorf("getting text from update: text not found")
	}
	return r.re.MatchString(text), nil
}

func TextMatchesRegexp(pattern string) core.Condition {
	return textMatchesRegexpCondition{
		re: regexp.MustCompile(pattern),
	}
}

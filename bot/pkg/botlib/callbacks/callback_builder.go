package callbacks

import (
	"fmt"
	"strings"

	"github.com/andreychh/coopera-bot/pkg/immutable"
)

type callbackBuilder struct {
	prefix string
	params immutable.Map[string, string]
}

func (c callbackBuilder) Encode() string {
	parts := make([]string, 0, c.params.Len()+1)
	parts = append(parts, c.prefix)
	for key, val := range c.params.All() {
		parts = append(parts, fmt.Sprintf("%s=%s", key, val))
	}
	return strings.Join(parts, ":")
}

func (c callbackBuilder) With(key string, value string) CallbackData {
	return callbackBuilder{
		prefix: c.prefix,
		params: c.params.With(key, value),
	}
}

func Builder(prefix string) CallbackData {
	return callbackBuilder{
		prefix: prefix,
		params: immutable.EmptyMap[string, string](),
	}
}

package callbacks

import (
	"fmt"
	"strings"

	"github.com/andreychh/coopera-bot/pkg/immutable"
)

type oncomingData struct {
	prefix string
	params immutable.Map[string, string]
}

func (o oncomingData) With(key string, value string) Outgoing {
	return oncomingData{
		prefix: o.prefix,
		params: o.params.With(key, value),
	}
}

func (o oncomingData) String() string {
	parts := make([]string, 0, o.params.Len()+1)
	parts = append(parts, o.prefix)
	for key, val := range o.params.All() {
		parts = append(parts, fmt.Sprintf("%s=%s", key, val))
	}
	return strings.Join(parts, ":")
}

func OutcomingData(prefix string) Outgoing {
	return oncomingData{
		prefix: prefix,
		params: immutable.EmptyMap[string, string](),
	}
}

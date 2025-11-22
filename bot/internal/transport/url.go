package transport

import (
	"fmt"
	"net/url"

	"github.com/andreychh/coopera-bot/pkg/immutable"
)

type OutcomingURL interface {
	With(key, value string) OutcomingURL
	String() string
}

type outcomingURL struct {
	url    string
	values immutable.Map[string, string]
}

func (o outcomingURL) With(key, value string) OutcomingURL {
	return outcomingURL{
		url:    o.url,
		values: o.values.With(key, value),
	}
}

func (o outcomingURL) String() string {
	if o.values.Len() == 0 {
		return o.url
	}
	values := url.Values{}
	for key, value := range o.values.All() {
		values.Add(key, value)
	}
	return fmt.Sprintf("%s?%s", o.url, values.Encode())
}

func NewOutcomingURL(url string) OutcomingURL {
	return outcomingURL{
		url:    url,
		values: immutable.EmptyMap[string, string](),
	}
}

package buttons

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type urlButton struct {
	text string
	url  string
}

func (u urlButton) AsObject() repr.Object {
	return json.Object().
		WithField("text", json.String(u.text)).
		WithField("url", json.String(u.url))
}

func URLButton(text string, url string) InlineButton {
	return urlButton{
		text: text,
		url:  url,
	}
}

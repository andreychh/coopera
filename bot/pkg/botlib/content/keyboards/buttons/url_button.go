package buttons

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type urlButton struct {
	text string
	url  string
}

func (u urlButton) Structure() repr.Structure {
	return json.Object(json.Fields{
		"text": json.Str(u.text),
		"url":  json.Str(u.url),
	})
}

func URLButton(text string, url string) InlineButton {
	return urlButton{
		text: text,
		url:  url,
	}
}

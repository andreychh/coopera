package buttons

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type textButton struct {
	text string
}

func (t textButton) AsObject() repr.Object {
	return json.Object().
		WithField("text", json.String(t.text))
}

func TextButton(text string) ReplyButton {
	return textButton{
		text: text,
	}
}

package content

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type textContent struct {
	text string
}

func (t textContent) AsObject() repr.Object {
	return json.Object().WithField("text", json.String(t.text))
}

func Text(text string) ObjectContent {
	return textContent{text: text}
}

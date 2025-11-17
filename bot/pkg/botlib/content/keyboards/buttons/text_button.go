package buttons

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type textButton struct {
	text string
}

func (t textButton) Structure() repr.Structure {
	return json.Object(json.Fields{
		"text": json.Str(t.text),
	})
}

func TextButton(text string) ReplyButton {
	return textButton{
		text: text,
	}
}

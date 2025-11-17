package buttons

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type callbackButton struct {
	text         string
	callbackData string
}

func (c callbackButton) Structure() repr.Structure {
	return json.Object(json.Fields{
		"text":          json.Str(c.text),
		"callback_data": json.Str(c.callbackData),
	})
}

func CallbackButton(text string, callbackData string) InlineButton {
	return callbackButton{
		text:         text,
		callbackData: callbackData,
	}
}

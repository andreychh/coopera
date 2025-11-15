package buttons

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type callbackButton struct {
	text         string
	callbackData string
}

func (c callbackButton) AsObject() repr.Object {
	return json.Object().
		WithField("text", json.String(c.text)).
		WithField("callback_data", json.String(c.callbackData))
}

func CallbackButton(text string, callbackData string) InlineButton {
	return callbackButton{
		text:         text,
		callbackData: callbackData,
	}
}

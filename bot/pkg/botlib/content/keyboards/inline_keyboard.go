package keyboards

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type inlineKeyboard struct {
	origin  content.ObjectContent
	buttons buttons.ButtonMatrix
}

func (i inlineKeyboard) AsObject() repr.Object {
	return i.origin.AsObject().WithField(
		"reply_markup",
		json.Object().WithField("inline_keyboard", i.buttons.AsArray()),
	)
}

func InlineKeyboard(origin content.ObjectContent, buttons buttons.ButtonMatrix) Keyboard {
	return inlineKeyboard{
		origin:  origin,
		buttons: buttons,
	}
}

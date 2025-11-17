package keyboards

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type inlineKeyboard struct {
	origin  content.Content
	buttons buttons.ButtonMatrix[buttons.InlineButton]
}

func (i inlineKeyboard) Structure() repr.Structure {
	return repr.Must(i.origin.Structure().Extend(
		repr.EmptyPath(),
		json.Object(json.Fields{
			"reply_markup": json.Object(json.Fields{
				"inline_keyboard": i.buttons.Structure(),
			}),
		}),
	))
}

func (i inlineKeyboard) Method() string {
	return i.origin.Method()
}

func Inline(content content.Content, buttons buttons.ButtonMatrix[buttons.InlineButton]) InlineKeyboard {
	return inlineKeyboard{
		origin:  content,
		buttons: buttons,
	}
}

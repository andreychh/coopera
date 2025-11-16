package keyboards

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type emptyKeyboard struct {
	origin content.Content
}

func (e emptyKeyboard) Structure() repr.Structure {
	return repr.Must(e.origin.Structure().Extend(
		repr.EmptyPath(),
		json.Object(json.Fields{
			"reply_markup": json.Object(json.Fields{
				"remove_keyboard": json.Bool(true),
			}),
		}),
	))
}

func (e emptyKeyboard) Method() string {
	return e.origin.Method()
}

func Empty(content content.Content) ReplyKeyboard {
	return emptyKeyboard{
		origin: content,
	}
}

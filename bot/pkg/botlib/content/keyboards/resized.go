package keyboards

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type resized struct {
	origin ReplyKeyboard
}

func (r resized) Structure() repr.Structure {
	return repr.Must(r.origin.Structure().Extend(
		repr.PathOf("reply_markup"),
		json.Object(json.Fields{
			"resize_keyboard": json.Bool(true),
		}),
	))
}

func (r resized) Method() string {
	return r.origin.Method()
}

func Resized(keyboard ReplyKeyboard) ReplyKeyboard {
	return resized{origin: keyboard}
}

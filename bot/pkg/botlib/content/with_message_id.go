package content

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type withMessageID struct {
	origin Content
	id     int
}

func (w withMessageID) Structure() repr.Structure {
	return repr.Must(w.origin.Structure().Extend(
		repr.EmptyPath(),
		json.Object(json.Fields{
			"message_id": json.Int(w.id),
		}),
	))
}

func (w withMessageID) Method() string {
	return w.origin.Method()
}

func WithMessageID(content Content, id int) Content {
	return withMessageID{
		origin: content,
		id:     id,
	}
}

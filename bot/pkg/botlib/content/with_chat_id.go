package content

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type withChatID struct {
	origin Content
	id     int64
}

func (w withChatID) Structure() repr.Structure {
	return repr.Must(w.origin.Structure().Extend(
		repr.EmptyPath(),
		json.Object(json.Fields{
			"chat_id": json.I64(w.id),
		}),
	))
}

func (w withChatID) Method() string {
	return w.origin.Method()
}

func WithChatID(content Content, recipientID int64) Content {
	return withChatID{
		origin: content,
		id:     recipientID,
	}
}

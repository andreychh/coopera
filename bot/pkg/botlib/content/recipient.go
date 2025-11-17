package content

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type Recipient struct {
	origin Content
	id     int64
}

func (r Recipient) Structure() repr.Structure {
	return repr.Must(r.origin.Structure().Extend(
		repr.EmptyPath(),
		json.Object(json.Fields{
			"chat_id": json.I64(r.id),
		}),
	))
}

func (r Recipient) Method() string {
	return r.origin.Method()
}

func To(content Content, recipientID int64) Content {
	return Recipient{
		origin: content,
		id:     recipientID,
	}
}

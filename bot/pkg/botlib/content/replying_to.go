package content

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type replyingMessage struct {
	id      int64
	content Content
}

func (r replyingMessage) Structure() repr.Structure {
	return repr.Must(r.content.Structure().Extend(
		repr.EmptyPath(),
		json.Object(json.Fields{
			"reply_to_message_id": json.I64(r.id),
		}),
	))
}

func (r replyingMessage) Method() string {
	return r.content.Method()
}

func ReplyingTo(messageID int64, content Content) Content {
	return replyingMessage{
		id:      messageID,
		content: content,
	}
}

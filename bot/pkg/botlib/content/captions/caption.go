package captions

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type caption struct {
	origin content.Content
	text   string
}

func (c caption) Structure() repr.Structure {
	return repr.Must(c.origin.Structure().Extend(
		repr.EmptyPath(),
		json.Object(json.Fields{
			"caption": json.Str(c.text),
		}),
	))
}

func (c caption) Method() string {
	return c.origin.Method()
}

func Caption(content content.Content, text string) content.Content {
	return caption{origin: content, text: text}
}

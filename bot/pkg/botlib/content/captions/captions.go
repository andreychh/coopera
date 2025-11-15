package captions

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type Caption interface {
	AsObject() repr.Object
}

type caption struct {
	origin content.ObjectContent
	text   string
}

func (c caption) AsObject() repr.Object {
	return c.origin.AsObject().WithField("caption", json.String(c.text))
}

func New(text string) content.ObjectContent {
	return caption{text: text}
}

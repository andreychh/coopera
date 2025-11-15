package content

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type urlPhotoContent struct {
	url string
}

func (u urlPhotoContent) AsObject() repr.Object {
	return json.Object().WithField("photo", json.String(u.url))
}

func URLPhoto(url string) ObjectContent {
	return urlPhotoContent{url: url}
}

package content

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
)

type urlPhotoContent struct {
	url string
}

func (u urlPhotoContent) Structure() repr.Structure {
	return json.Object(json.Fields{
		"photo": json.Str(u.url),
	})
}

func (u urlPhotoContent) Method() string {
	return "sendPhoto"
}

func URLPhoto(url string) Content {
	return urlPhotoContent{url: url}
}

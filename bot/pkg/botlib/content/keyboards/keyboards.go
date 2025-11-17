package keyboards

import "github.com/andreychh/coopera-bot/pkg/botlib/content"

type Keyboard interface {
	content.Content
}

type InlineKeyboard interface {
	Keyboard
}

type ReplyKeyboard interface {
	Keyboard
}

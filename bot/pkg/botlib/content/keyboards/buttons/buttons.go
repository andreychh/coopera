package buttons

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
)

type ButtonMatrix[T Button] interface {
	WithRow(row ButtonRow[T]) ButtonMatrix[T]
	Structure() repr.Structure
}

type ButtonRow[T Button] interface {
	WithButton(button T) ButtonRow[T]
	Structure() repr.Structure
}

type Button interface {
	Structure() repr.Structure
}

type InlineButton interface {
	Button
}

type ReplyButton interface {
	Button
}

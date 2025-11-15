package buttons

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
)

type ButtonMatrix interface {
	WithRow(row ButtonRow) ButtonMatrix
	AsArray() repr.Array
}

type ButtonRow interface {
	WithButton(button Button) ButtonRow
	AsArray() repr.Array
}

type Button interface {
	AsObject() repr.Object
}

type InlineButton interface {
	Button
}

type ReplyButton interface {
	Button
}

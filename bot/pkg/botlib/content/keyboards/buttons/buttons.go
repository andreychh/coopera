package buttons

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
)

type ButtonMatrix interface {
	WithRow(row ButtonRow) ButtonMatrix
	Structure() repr.Structure
}

type ButtonRow interface {
	WithButton(button Button) ButtonRow
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

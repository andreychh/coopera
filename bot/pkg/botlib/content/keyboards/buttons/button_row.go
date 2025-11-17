package buttons

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
	"github.com/andreychh/coopera-bot/pkg/slices"
)

type buttonRow[T Button] struct {
	buttons []T
}

func (r buttonRow[T]) WithButton(button T) ButtonRow[T] {
	return buttonRow[T]{buttons: slices.With(r.buttons, button)}
}

func (r buttonRow[T]) Structure() repr.Structure {
	elements := make([]repr.Structure, len(r.buttons))
	for i, btn := range r.buttons {
		elements[i] = btn.Structure()
	}
	return json.Array(elements...)
}

func Row[T Button](buttons ...T) ButtonRow[T] {
	return buttonRow[T]{buttons: buttons}
}

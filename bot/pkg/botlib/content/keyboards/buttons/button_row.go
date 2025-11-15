package buttons

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
	"github.com/andreychh/coopera-bot/pkg/slices"
)

type buttonRow struct {
	buttons []Button
}

func (r buttonRow) WithButton(button Button) ButtonRow {
	return buttonRow{buttons: slices.With(r.buttons, button)}
}

func (r buttonRow) AsArray() repr.Array {
	array := json.Array()
	for _, btn := range r.buttons {
		array = array.WithElement(btn.AsObject())
	}
	return array
}

func Row(buttons ...Button) ButtonRow {
	return buttonRow{buttons: buttons}
}

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

func (r buttonRow) Structure() repr.Structure {
	elements := make([]repr.Structure, len(r.buttons))
	for i, btn := range r.buttons {
		elements[i] = btn.Structure()
	}
	return json.Array(elements...)
}

func Row(buttons ...Button) ButtonRow {
	return buttonRow{buttons: buttons}
}

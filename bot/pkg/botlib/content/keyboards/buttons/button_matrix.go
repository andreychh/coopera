package buttons

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
	"github.com/andreychh/coopera-bot/pkg/slices"
)

type buttonMatrix struct {
	Rows []ButtonRow
}

func (m buttonMatrix) WithRow(row ButtonRow) ButtonMatrix {
	return buttonMatrix{Rows: slices.With(m.Rows, row)}
}

func (m buttonMatrix) AsArray() repr.Array {
	array := json.Array()
	for _, row := range m.Rows {
		array = array.WithElement(row.AsArray())
	}
	return array
}

func Matrix(rows ...ButtonRow) ButtonMatrix {
	return buttonMatrix{Rows: rows}
}

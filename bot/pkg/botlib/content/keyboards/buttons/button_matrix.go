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

func (m buttonMatrix) Structure() repr.Structure {
	elements := make([]repr.Structure, len(m.Rows))
	for i, row := range m.Rows {
		elements[i] = row.Structure()
	}
	return json.Array(elements...)
}

func Matrix(rows ...ButtonRow) ButtonMatrix {
	return buttonMatrix{Rows: rows}
}

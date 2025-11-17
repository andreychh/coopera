package buttons

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
	"github.com/andreychh/coopera-bot/pkg/slices"
)

type buttonMatrix[T Button] struct {
	rows []ButtonRow[T]
}

func (m buttonMatrix[T]) WithRow(row ButtonRow[T]) ButtonMatrix[T] {
	return buttonMatrix[T]{rows: slices.With(m.rows, row)}
}

func (m buttonMatrix[T]) Structure() repr.Structure {
	elements := make([]repr.Structure, len(m.rows))
	for i, row := range m.rows {
		elements[i] = row.Structure()
	}
	return json.Array(elements...)
}

func Matrix[T Button](rows ...ButtonRow[T]) ButtonMatrix[T] {
	return buttonMatrix[T]{rows: rows}
}

package buttons

import (
	"github.com/andreychh/coopera-bot/pkg/repr"
	"github.com/andreychh/coopera-bot/pkg/repr/json"
	"github.com/andreychh/coopera-bot/pkg/slices"
)

type buttonMatrix[T Button] struct {
	Rows []ButtonRow[T]
}

func (m buttonMatrix[T]) WithRow(row ButtonRow[T]) ButtonMatrix[T] {
	return buttonMatrix[T]{Rows: slices.With(m.Rows, row)}
}

func (m buttonMatrix[T]) Structure() repr.Structure {
	elements := make([]repr.Structure, len(m.Rows))
	for i, row := range m.Rows {
		elements[i] = row.Structure()
	}
	return json.Array(elements...)
}

func Matrix[T Button](rows ...ButtonRow[T]) ButtonMatrix[T] {
	return buttonMatrix[T]{Rows: rows}
}

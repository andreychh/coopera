package hsm

import (
	"iter"

	"github.com/andreychh/coopera-bot/pkg/immutable"
)

type path struct {
	states immutable.Slice[State]
}

func (p path) All() iter.Seq[State] {
	return func(yield func(State) bool) {
		for _, s := range p.states.All() {
			if !yield(s) {
				return
			}
		}
	}
}

func (p path) Add(state State) path {
	return path{
		states: p.states.Insert(p.states.Len(), state),
	}
}

func (p path) Reversed() path {
	return path{
		states: immutable.ReversedSlice(p.states),
	}
}

func EmptyPath() path {
	return path{
		states: immutable.EmptySlice[State](),
	}
}

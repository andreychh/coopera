package hsm

import (
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
)

type Spec interface {
	ID() core.ID
	Compile(parent State) (State, []State, error)
}

type nodeSpec struct {
	id       core.ID
	behavior Behavior
	children Children
}

func (n nodeSpec) ID() core.ID {
	return n.id
}

func Node(id core.ID, b Behavior, kids Children) Spec {
	return nodeSpec{
		id:       id,
		behavior: b,
		children: kids,
	}
}

func (n nodeSpec) Compile(parent State) (State, []State, error) {
	st := &generalState{
		id:       n.id,
		behavior: n.behavior,
		parent:   parent,
		initial:  EdgeState(),
	}
	flatList := []State{st}
	var runtimeInitial State
	foundInitial := false
	for _, childSpec := range n.children.All() {
		childState, childDescendants, err := childSpec.Compile(st)
		if err != nil {
			return nil, nil, fmt.Errorf("node %q compile failed on child %q: %w", n.id, childSpec.ID(), err)
		}
		if childSpec.ID() == n.children.Initial().ID() {
			runtimeInitial = childState
			foundInitial = true
		}
		flatList = append(flatList, childDescendants...)
	}
	if !foundInitial {
		return nil, nil, fmt.Errorf("node %q invariant violation: initial child %q defined but not found in children list", n.id, n.children.Initial().ID())
	}
	st.initial = runtimeInitial
	return st, flatList, nil
}

type leafSpec struct {
	id       core.ID
	behavior Behavior
}

func (l leafSpec) ID() core.ID {
	return l.id
}

func (l leafSpec) Compile(parent State) (State, []State, error) {
	st := &generalState{
		id:       l.id,
		behavior: l.behavior,
		parent:   parent,
		initial:  EdgeState(),
	}
	return st, []State{st}, nil
}

func Leaf(id core.ID, b Behavior) Spec {
	return leafSpec{
		id:       id,
		behavior: b,
	}
}

func Anchor(id core.ID) Spec {
	return leafSpec{
		id:       id,
		behavior: NoBehavior(),
	}
}

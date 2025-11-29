package hsm

import (
	"fmt"
	"iter"
)

type cursor struct {
	node State
}

func (c cursor) Up() iter.Seq[State] {
	return func(yield func(State) bool) {
		curr := c.node
		for {
			if !yield(curr) {
				return
			}
			parent, ok := curr.Parent()
			if !ok {
				return
			}
			curr = parent
		}
	}
}

func (c cursor) Down() iter.Seq[State] {
	return func(yield func(State) bool) {
		curr := c.node
		for {
			if !yield(curr) {
				return
			}
			child, ok := curr.Initial()
			if !ok {
				return
			}
			curr = child
		}
	}
}

func (c cursor) Root() State {
	var root State
	for state := range c.Up() {
		root = state
	}
	return root
}

func (c cursor) Leaf() State {
	var leaf State
	for state := range c.Down() {
		leaf = state
	}
	return leaf
}

func (c cursor) UpPath(limit State) (path, error) {
	p := EmptyPath()
	for state := range c.Up() {
		if state.ID() == limit.ID() {
			return p, nil
		}
		p = p.Add(state)
	}
	return path{}, fmt.Errorf("limit state %q not found in path", limit.ID())
}

func (c cursor) DownPath(limit State) (path, error) {
	p := EmptyPath()
	for state := range c.Down() {
		if state.ID() == limit.ID() {
			return p, nil
		}
		p = p.Add(state)
	}
	return path{}, fmt.Errorf("limit state %q not found in path", limit.ID())
}

func Cursor(node State) cursor {
	return cursor{node: node}
}

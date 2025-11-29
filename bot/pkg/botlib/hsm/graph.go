package hsm

import (
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/immutable"
)

type Graph interface {
	State(id core.ID) (State, bool)
	Root() State
	Transition(fromID, nextID core.ID) (Transition, error)
}
type graph struct {
	root  State
	index immutable.Map[core.ID, State]
}

func (g graph) Root() State {
	return g.root
}

func (g graph) State(id core.ID) (State, bool) {
	return g.index.Get(id)
}

func (g graph) Transition(fromID, nextID core.ID) (Transition, error) {
	from, ok := g.State(fromID)
	if !ok {
		return Transition{}, fmt.Errorf("source state %q not found", fromID)
	}
	targetNode, ok := g.State(nextID)
	if !ok {
		return Transition{}, fmt.Errorf("target state %q not found", nextID)
	}
	realTarget := Cursor(targetNode).Leaf()
	lca := g.lca(from, realTarget)
	exitPath, err := Cursor(from).UpPath(lca)
	if err != nil {
		// Этого быть не должно, если граф корректен
		return Transition{}, fmt.Errorf("invariant: exit path not found: %w", err)
	}
	enterPathRaw, err := Cursor(realTarget).UpPath(lca)
	if err != nil {
		return Transition{}, fmt.Errorf("invariant: enter path not found: %w", err)
	}
	enterPath := enterPathRaw.Reversed()
	return Transition{
		targetID: realTarget.ID(),
		exit:     exitPath,
		enter:    enterPath,
	}, nil
}

func (g graph) lca(from, to State) State {
	if from.ID() == to.ID() {
		p, ok := from.Parent()
		if !ok {
			// Если корень делает переход сам в себя - LCA это сам корень (редкий кейс)
			return from
		}
		return p
	}
	ancestors := make(map[core.ID]struct{})
	for state := range Cursor(from).Up() {
		ancestors[state.ID()] = struct{}{}
	}
	for state := range Cursor(to).Up() {
		if _, exists := ancestors[state.ID()]; exists {
			return state
		}
	}
	panic(fmt.Sprintf("invariant violated: no common ancestor found between %q and %q", from.ID(), to.ID()))
}

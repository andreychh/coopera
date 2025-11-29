package hsm

import (
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/immutable"
)

type compiler struct {
	root Spec
}

func (c compiler) Graph() (Graph, error) {
	rootState, flatList, err := c.root.Compile(EdgeState())
	if err != nil {
		return nil, fmt.Errorf("compilation failed: %w", err)
	}
	index := immutable.EmptyMap[core.ID, State]()
	for _, s := range flatList {
		id := s.ID()
		_, exists := index.Get(id)
		if exists {
			return nil, fmt.Errorf("hsm compile: duplicate state ID %q detected", id)
		}
		index = index.With(id, s)
	}
	return &graph{
		root:  rootState,
		index: index,
	}, nil
}

func NewCompiler(root Spec) Compiler {
	return &compiler{root: root}
}

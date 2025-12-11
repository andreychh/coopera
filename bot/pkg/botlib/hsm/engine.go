package hsm

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/andreychh/coopera-bot/pkg/botlib/sessions"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attrs"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type engine struct {
	sessions sessions.Sessions
	graph    Graph
}

func (e engine) TryExecute(ctx context.Context, update telegram.Update) (bool, error) {
	id, ok := attrs.ChatID(update).Value()
	if !ok {
		return false, nil
	}
	currState, err := e.loadAndValidateState(ctx, id, update)
	if err != nil {
		return false, err
	}
	decision, err := Chain(Cursor(currState).Up()).Handle(ctx, update)
	if err != nil {
		return false, fmt.Errorf("handling error: %w", err)
	}
	if !decision.HasTransition() {
		return decision.Handled(), nil
	}
	route, err := e.graph.Transition(currState.ID(), decision.Next())
	if err != nil {
		return false, fmt.Errorf("transition calculation failed: %w", err)
	}
	if err := route.PerformCleanup(ctx, update); err != nil {
		slog.Warn("cleanup phase error", "error", err)
	}
	if err := e.sessions.Session(id).ChangeStateID(ctx, route.TargetID()); err != nil {
		return false, fmt.Errorf("commit failed: %w", err)
	}
	if err := route.PerformPrompt(ctx, update); err != nil {
		return false, fmt.Errorf("prompt phase failed: %w", err)
	}
	return true, nil
}

func (e engine) loadAndValidateState(ctx context.Context, chatID int64, u telegram.Update) (State, error) {
	currID, err := e.sessions.Session(chatID).StateID(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting session: %w", err)
	}
	state, exists := e.graph.State(currID)
	if currID == "" || !exists {
		return e.emergencyReset(ctx, chatID, u)
	}
	if _, isNode := state.Initial(); isNode {
		realLeaf := Cursor(state).Leaf()
		for s := range Cursor(state).Down() {
			if s.ID() == state.ID() {
				continue
			}
			if err := s.Enter(ctx, u); err != nil {
				slog.Error("floating recovery failed", "state", s.ID(), "error", err)
				return e.emergencyReset(ctx, chatID, u)
			}
		}
		if err := e.sessions.Session(chatID).ChangeStateID(ctx, realLeaf.ID()); err != nil {
			return nil, fmt.Errorf("floating fix commit failed: %w", err)
		}
		return realLeaf, nil
	}
	return state, nil
}

func (e engine) emergencyReset(ctx context.Context, chatID int64, u telegram.Update) (State, error) {
	rootLeaf := Cursor(e.graph.Root()).Leaf()
	var initPath []State
	for s := range Cursor(rootLeaf).Up() {
		initPath = append(initPath, s)
	}
	for i, j := 0, len(initPath)-1; i < j; i, j = i+1, j-1 {
		initPath[i], initPath[j] = initPath[j], initPath[i]
	}
	for _, s := range initPath {
		if err := s.Enter(ctx, u); err != nil {
			return nil, fmt.Errorf("cold start enter failed at %q: %w", s.ID(), err)
		}
	}
	if err := e.sessions.Session(chatID).ChangeStateID(ctx, rootLeaf.ID()); err != nil {
		return nil, fmt.Errorf("cold start commit failed: %w", err)
	}
	return rootLeaf, nil
}

func NewEngine(s sessions.Sessions, g Graph) *engine {
	return &engine{sessions: s, graph: g}
}

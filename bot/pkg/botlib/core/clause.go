package core

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Clause represents a single processing step in the update pipeline.
//
// A Clause is the fundamental building block of the routing tree, responsible
// for coordinating Condition checks and Action execution.
//
// Return values:
//   - (true, nil): Update was handled successfully. Abort further clauses.
//   - (false, nil): Clause was not applicable. Try the next clause.
//   - (_, err): Error occurred during execution. Abort chain.
//
// Common wrappers:
//   - routing.FirstExecuted(...): runs clauses until one returns (true, nil).
//   - routing.If(...): wraps a Clause with a pre-check Condition.
//
// # Principle of Tree Composition
//
// The entire routing and processing logic MUST be composed exclusively of
// objects that implement the Clause interface. This ensures uniformity,
// polymorphism, and a robust hierarchical structure (Composite Pattern).
type Clause interface {
	TryExecute(ctx context.Context, update telegram.Update) (executed bool, err error)
}

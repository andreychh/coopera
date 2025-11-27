package core

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Condition defines a pure predicate on an Update.
//
// Conditions MUST NOT cause any external system state mutation (writes to DB,
// sending messages). They are strictly used for evaluation and reading state,
// and must be composable via And/Or/Not.
//
// Return values:
//   - (true, nil): The condition holds (business rule met).
//   - (false, nil): The condition failed (business rule not met).
//   - (_, err): Failure during evaluation (structural/I/O error).
//
// # Principle of Failure Separation
//
// A Condition MUST return (false, nil) if the business rule is not met (e.g.,
// Command is "help" but expected "start"). This is a **normal execution path**
// that allows routing to continue.
//
// A Condition MUST return (_, error) if there is a structural or I/O failure
// (e.g., database lookup failed, or the Update structurally LACKS the required
// data for evaluation). This is an **exceptional path** that signals the router
// to abort the chain.
type Condition interface {
	Holds(ctx context.Context, update telegram.Update) (held bool, err error)
}

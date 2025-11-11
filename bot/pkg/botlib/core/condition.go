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
//
// Conditions are categorized into two primary types based on their role in the
// routing tree:
//   - Guards (e.g., HasCommand, HasText):
//     These are fast, simple conditions that check for structural data presence.
//     They MUST always return (false, nil) if the data is missing, guaranteeing
//     safety.
//   - Business Conditions (e.g., CommandIn, DialogueExists):
//     These check specific application state or business rules.
//
// The implementation provides pairs of constructors for Business Conditions:
//   - Unsecured (e.g., CommandIn):
//     These are simple wrappers that WILL return a structural error if required
//     data (like ChatID or Command) is missing. They are designed for strict
//     routing paths where data is guaranteed.
//   - Secured (e.g., SafeCommandIn):
//     These automatically compose the necessary Guard (e.g., HasCommand()),
//     ensuring the Condition returns (false, nil) when structural data is missing,
//     preventing unnecessary routing interruptions.
type Condition interface {
	Holds(ctx context.Context, update telegram.Update) (held bool, err error)
}

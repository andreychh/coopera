package core

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Action represents the final, executable operation in a Clause.
//
// Actions are the sole source of state mutation and side effects (writes to DB,
// sending messages, I/O that changes external state).
//
// # Principle of Pure Mutation
//
// An Action MUST focus only on performing its designated state mutation (I/O).
//
// Return values:
//   - nil: Action completed successfully.
//   - err: Failure to perform the action (I/O error, database error).
type Action interface {
	Perform(ctx context.Context, update telegram.Update) error
}

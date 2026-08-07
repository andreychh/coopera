// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import "fmt"

// UserNotFoundError says the actor of a request is nobody the system
// knows. While X-User-Id stands in for authentication, each usecase
// discovers this for itself; real authentication settles it before a
// usecase runs, and this error leaves the usecases' failure sets.
type UserNotFoundError struct {
	ID ID
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user %s not found", e.ID)
}

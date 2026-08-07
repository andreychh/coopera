// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// UnexpectedError is a failure the domain has no answer for: the
// database being unreachable, a query breaking, a row the schema should
// have ruled out. It carries the cause for logging, and callers turn it
// into a 500 without inspecting it.
type UnexpectedError struct {
	Err error
}

func (e UnexpectedError) Error() string {
	return e.Err.Error()
}

// Unwrap exposes the cause to [errors.Is] and [errors.As].
func (e UnexpectedError) Unwrap() error {
	return e.Err
}

func (UnexpectedError) createInviteLinkError() {}

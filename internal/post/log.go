// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package post carries sign-in codes to the addresses that asked for
// them.
package post

import (
	"context"
	"log"

	"github.com/andreychh/coopera/internal/domain"
)

// Log writes the code where a developer can read it instead of mailing
// it anywhere. It exists so that signing in works before there is any
// mail to send it with, and it must never run outside development: a
// code printed to a log is a code anyone with the log can use.
type Log struct{}

func NewLog() Log {
	return Log{}
}

func (Log) Send(_ context.Context, to domain.Email, code domain.SignInCode) error {
	log.Printf("sign-in code for %s: %s", to, code)
	return nil
}

package dberrors

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("record not found")
	ErrAlreadyExists       = errors.New("record already exists")
	ErrDB                  = errors.New("database error")
	ErrTransactionNotFound = errors.New("transaction not found in context")
	ErrFailToCastScan      = errors.New("failed to cast scan result")
	ErrFailToAdd           = errors.New("failed to add member to team")
	ErrFailDelete          = errors.New("failed to delete team")
	ErrFailCreate          = errors.New("failed to create team")
	ErrFailGet             = errors.New("failed to get members")
	ErrFailCheckExists     = errors.New("failed to check existence")
)

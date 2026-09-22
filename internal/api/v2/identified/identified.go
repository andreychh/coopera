// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package identified gives each request a server-generated identifier.
package identified

import (
	"net/http"

	"github.com/andreychh/coopera/internal/api/v2/raw"
	"github.com/google/uuid"
)

// Request is a [raw.Request] with an identifier assigned at a [Checkpoint].
//
// A Request is obtained only from a Checkpoint, so the identifier is always
// server-generated.
type Request struct {
	raw.Request

	id uuid.UUID
}

// ID returns the identifier the request was given.
func (r Request) ID() uuid.UUID {
	return r.id
}

// Handler responds to a request that has an identifier.
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r Request)
}

// HandlerFunc adapts an ordinary function to a [Handler].
type HandlerFunc func(w http.ResponseWriter, r Request)

// ServeHTTP implements [Handler] by calling f.
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r Request) {
	f(w, r)
}

// Checkpoint is a [raw.Handler] that gives each request a new identifier and
// reports it in the X-Request-ID response header. An X-Request-ID sent by the
// client is never used as the identifier.
//
// The zero Checkpoint panics when it serves.
type Checkpoint struct {
	next Handler
}

// NewCheckpoint constructs a Checkpoint that passes identified requests to
// next.
func NewCheckpoint(next Handler) Checkpoint {
	return Checkpoint{next: next}
}

// ServeHTTP implements [raw.Handler]. It panics only if a random source
// installed with [uuid.SetRand] fails.
func (c Checkpoint) ServeHTTP(w http.ResponseWriter, r raw.Request) {
	id := uuid.Must(uuid.NewV7())
	w.Header().Set("X-Request-ID", id.String())
	c.next.ServeHTTP(w, Request{Request: r, id: id})
}

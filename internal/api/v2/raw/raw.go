// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package raw narrows an incoming HTTP request to the parts a handler may
// read.
package raw

import (
	"context"
	"io"
	"net/http"
)

// Request is an incoming HTTP request narrowed to what a handler may read.
// Nothing it carries has been checked.
//
// A Request is obtained only from a [Checkpoint]; the zero Request panics on
// every method.
type Request struct {
	origin *http.Request
}

// Context returns the request's context.
func (r Request) Context() context.Context {
	return r.origin.Context()
}

// Header returns the first value of the named header, or "" if the request
// does not carry it. The name is matched without regard to case.
func (r Request) Header(name string) string {
	return r.origin.Header.Get(name)
}

// PathValue returns the value of the named wildcard in the route pattern the
// request matched, or "" if the pattern has no such wildcard.
func (r Request) PathValue(name string) string {
	return r.origin.PathValue(name)
}

// Query returns the first value of the named query parameter and reports
// whether the request carries the parameter. A parameter present without a
// value yields "" and true.
func (r Request) Query(name string) (string, bool) {
	values := r.origin.URL.Query()
	return values.Get(name), values.Has(name)
}

// Body returns the request body. It is never nil: a request without a body
// yields [io.EOF] at once.
func (r Request) Body() io.Reader {
	return r.origin.Body
}

// Handler responds to a [Request].
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r Request)
}

// HandlerFunc adapts an ordinary function to a [Handler].
type HandlerFunc func(w http.ResponseWriter, r Request)

// ServeHTTP implements [Handler] by calling f.
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r Request) {
	f(w, r)
}

// Checkpoint is an [http.Handler] that passes each request on as a [Request].
//
// The zero Checkpoint panics when it serves.
type Checkpoint struct {
	next Handler
}

// NewCheckpoint constructs a Checkpoint that passes narrowed requests to next.
func NewCheckpoint(next Handler) Checkpoint {
	return Checkpoint{next: next}
}

// ServeHTTP implements [http.Handler].
func (c Checkpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.next.ServeHTTP(w, Request{origin: r})
}

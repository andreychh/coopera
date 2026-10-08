// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package raw narrows an incoming HTTP request to the parts a handler may
// read.
package raw

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/andreychh/coopera/internal/api/v2/pipeline"
)

// Request is an incoming HTTP request narrowed to what a handler may read.
// Nothing it carries has been checked.
//
// A Request is obtained only from a [Checkpoint]; the zero Request panics on
// every method.
type Request struct {
	origin *http.Request
}

// Context returns the request's context. It is canceled once the time the
// [Checkpoint] allows a request has passed, and when the client goes away.
func (r Request) Context() context.Context {
	return r.origin.Context()
}

// RemoteAddr returns the network address the request arrived from, as
// "IP:port" for a TCP connection. It is the address of the connection's peer:
// behind a proxy, it is the proxy's.
func (r Request) RemoteAddr() string {
	return r.origin.RemoteAddr
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

// Body returns the request body, capped at the size the [Checkpoint] allows.
// It is never nil: a request without a body yields [io.EOF] at once, and one
// longer than the cap yields a [*http.MaxBytesError].
func (r Request) Body() io.Reader {
	return r.origin.Body
}

// Handler responds to a [Request].
type Handler = pipeline.Handler[Request]

// HandlerFunc adapts an ordinary function to a [Handler].
type HandlerFunc = pipeline.HandlerFunc[Request]

// Checkpoint is an [http.Handler] that passes each request on as a [Request]
// under two bounds of the route it serves: how large a body it takes and how
// long it may take to answer. A handler reading past the body bound gets a
// [*http.MaxBytesError] and the connection is closed; past the time bound the
// request's context is canceled. Bounds belong to the route, so a route that
// takes uploads or runs long gets a Checkpoint of its own.
//
// Canceling the context stops only the work that watches it, and nothing
// answers on its own: the handler that finds its context canceled writes the
// response.
//
// The zero Checkpoint panics when it serves.
type Checkpoint struct {
	next    Handler
	body    int64
	timeout time.Duration
}

// NewCheckpoint constructs a Checkpoint that passes narrowed requests to next,
// caps their bodies at body bytes and gives each request timeout to be
// answered.
func NewCheckpoint(next Handler, body int64, timeout time.Duration) Checkpoint {
	return Checkpoint{next: next, body: body, timeout: timeout}
}

// ServeHTTP implements [http.Handler].
func (c Checkpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), c.timeout)
	defer cancel()

	r = r.WithContext(ctx)
	r.Body = http.MaxBytesReader(w, r.Body, c.body)
	c.next.ServeHTTP(w, Request{origin: r})
}

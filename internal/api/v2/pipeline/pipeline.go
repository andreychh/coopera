// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package pipeline declares what every element of the v2 request pipeline has
// in common.
package pipeline

import "net/http"

// Handler responds to a request of type R. A stage instantiates it with the
// request it passes on; an element that leaves the request as it found it is
// generic over R.
type Handler[R any] interface {
	ServeHTTP(w http.ResponseWriter, r R)
}

// HandlerFunc adapts an ordinary function to a [Handler].
type HandlerFunc[R any] func(w http.ResponseWriter, r R)

// ServeHTTP implements [Handler] by calling f.
func (f HandlerFunc[R]) ServeHTTP(w http.ResponseWriter, r R) {
	f(w, r)
}

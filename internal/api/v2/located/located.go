// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package located establishes the network address a request came from.
package located

import (
	"net/http"
	"net/netip"

	"github.com/andreychh/coopera/internal/api/v2/identified"
	"github.com/andreychh/coopera/internal/api/v2/pipeline"
)

// Request is an [identified.Request] with the address of the client that sent
// it.
//
// A Request is obtained only from a [Checkpoint].
type Request struct {
	identified.Request

	client netip.Addr
}

// ClientAddr returns the address of the client the request came from. It is
// always valid, and an IPv4 address is never in its IPv4-mapped IPv6 form.
func (r Request) ClientAddr() netip.Addr {
	return r.client
}

// ClientNetwork returns the network that identifies the client: its IPv4
// address as a /32, or the /64 its IPv6 address belongs to. An IPv6 client
// picks addresses from its /64 freely, so the address alone does not identify
// it.
func (r Request) ClientNetwork() netip.Prefix {
	bits := 64
	if r.client.Is4() {
		bits = 32
	}
	return netip.PrefixFrom(r.client, bits).Masked()
}

// Handler responds to a request whose client address is known.
type Handler = pipeline.Handler[Request]

// HandlerFunc adapts an ordinary function to a [Handler].
type HandlerFunc = pipeline.HandlerFunc[Request]

// Checkpoint is an [identified.Handler] that takes the client address from the
// connection the request arrived on. A request whose connection has no IP
// address and port, as on a Unix socket, is answered with 500 and goes no
// further: the address comes from the server, so the fault is the server's.
//
// The zero Checkpoint panics when it passes a request on.
type Checkpoint struct {
	next Handler
}

// NewCheckpoint constructs a Checkpoint that passes located requests to next.
func NewCheckpoint(next Handler) Checkpoint {
	return Checkpoint{next: next}
}

// ServeHTTP implements [identified.Handler].
func (c Checkpoint) ServeHTTP(w http.ResponseWriter, r identified.Request) {
	peer, err := netip.ParseAddrPort(r.RemoteAddr())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.next.ServeHTTP(w, Request{Request: r, client: peer.Addr().Unmap()})
}

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package limit refuses requests that come faster than a rate allows.
package limit

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/andreychh/coopera/internal/api/v2/pipeline"
	"github.com/andreychh/coopera/internal/ratelimiting"
)

// Limit is a [pipeline.Handler] that passes a request on only if it takes a token for
// it: from one bucket all requests share, or from the bucket of the request's
// key. A refused request is answered with the status given to the constructor
// and a Retry-After, the wait in whole seconds rounded up.
//
// The zero Limit panics when it serves.
type Limit[R any] struct {
	next   pipeline.Handler[R]
	take   func(r R, now time.Time) error
	status int
}

// Shared constructs a Limit that takes the token of every request from bucket,
// answers a refused request with status and passes the rest to next.
func Shared[R any](next pipeline.Handler[R], bucket *ratelimiting.Bucket, status int) Limit[R] {
	return Limit[R]{
		next: next,
		take: func(_ R, now time.Time) error {
			return bucket.Take(now, 1)
		},
		status: status,
	}
}

// Keyed constructs a Limit that keys each request with key, takes its token
// from the bucket of that key in buckets, answers a refused request with status
// and passes the rest to next.
func Keyed[R any, K comparable](
	next pipeline.Handler[R],
	buckets *ratelimiting.Buckets[K],
	key func(R) K,
	status int,
) Limit[R] {
	return Limit[R]{
		next: next,
		take: func(r R, now time.Time) error {
			return buckets.Take(key(r), now, 1)
		},
		status: status,
	}
}

// ServeHTTP implements [pipeline.Handler].
func (l Limit[R]) ServeHTTP(w http.ResponseWriter, r R) {
	err := l.take(r, time.Now())
	if err != nil {
		refuse(w, l.status, err)
		return
	}
	l.next.ServeHTTP(w, r)
}

// refuse answers with status and a Retry-After taken from err. An error that
// is not a [ratelimiting.LimitedError] means the bucket cannot serve any
// request, and is answered with 500 instead.
func refuse(w http.ResponseWriter, status int, err error) {
	limited, ok := errors.AsType[ratelimiting.LimitedError](err)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	seconds := (limited.Wait + time.Second - 1) / time.Second
	w.Header().Set("Retry-After", strconv.FormatInt(int64(seconds), 10))
	w.WriteHeader(status)
}

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ratelimiting_test

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/andreychh/coopera/internal/ratelimiting"
)

func TestNewBuckets(t *testing.T) {
	t.Parallel()

	t.Run("rejects a size below one", func(t *testing.T) {
		t.Parallel()

		for _, size := range []int{0, -1} {
			_, err := ratelimiting.NewBuckets[string](mustPolicy(t, 3), size)
			if err == nil {
				t.Errorf("NewBuckets(size %d) = nil error, want one", size)
			}
		}
	})
}

func TestBuckets_Take(t *testing.T) {
	t.Parallel()

	newBuckets := func(t *testing.T, size int) *ratelimiting.Buckets[string] {
		t.Helper()

		buckets, err := ratelimiting.NewBuckets[string](mustPolicy(t, 3), size)
		if err != nil {
			t.Fatalf("NewBuckets(size %d) = %v, want nil", size, err)
		}
		return buckets
	}

	// call is a Take made for key at start+at for n tokens, with the error it
	// must return.
	type call struct {
		key  string
		at   time.Duration
		n    int
		want error
	}
	const ms = time.Millisecond

	// Every case plays its calls in order on new buckets of capacity 3 that
	// earn one token per second, with room for two keys.
	cases := []struct {
		behavior string
		calls    []call
	}{{
		behavior: "keeps a separate bucket for each key",
		calls: []call{
			{key: "a", at: 0, n: 3, want: nil},
			{key: "a", at: 0, n: 1, want: limited(time.Second)},
			{key: "b", at: 0, n: 1, want: nil},
		},
	}, {
		behavior: "refuses a new key until a bucket refills",
		calls: []call{
			{key: "a", at: 0, n: 3, want: nil},
			{key: "b", at: 0, n: 3, want: nil},
			{key: "c", at: 0, n: 1, want: limited(3 * time.Second)},
		},
	}, {
		behavior: "serves known keys while it refuses new ones",
		calls: []call{
			{key: "a", at: 0, n: 3, want: nil},
			{key: "b", at: 0, n: 3, want: nil},
			{key: "c", at: 0, n: 1, want: limited(3 * time.Second)},
			{key: "a", at: time.Second, n: 1, want: nil},
		},
	}, {
		behavior: "makes room by dropping buckets that have refilled",
		calls: []call{
			{key: "a", at: 0, n: 3, want: nil},
			{key: "b", at: 0, n: 3, want: nil},
			{key: "c", at: 3 * time.Second, n: 3, want: nil},
			{key: "a", at: 3 * time.Second, n: 3, want: nil},
		},
	}, {
		behavior: "keeps a bucket that has not refilled",
		calls: []call{
			{key: "a", at: 0, n: 3, want: nil},
			{key: "b", at: 2 * time.Second, n: 3, want: nil},
			{key: "c", at: 2500 * ms, n: 1, want: limited(500 * ms)},
			{key: "a", at: 2500 * ms, n: 3, want: limited(500 * ms)},
		},
	}, {
		behavior: "counts buckets added since it last looked for room",
		calls: []call{
			{key: "a", at: 0, n: 3, want: nil},
			{key: "x", at: 0, n: 1, want: nil},
			{key: "b", at: time.Second, n: 1, want: nil},
			{key: "c", at: 2 * time.Second, n: 1, want: nil},
		},
	}, {
		behavior: "waits again when the bucket due to refill first is used",
		calls: []call{
			{key: "a", at: 0, n: 3, want: nil},
			{key: "b", at: time.Second, n: 3, want: nil},
			{key: "c", at: time.Second, n: 1, want: limited(2 * time.Second)},
			{key: "a", at: 2500 * ms, n: 1, want: nil},
			{key: "c", at: 3 * time.Second, n: 1, want: limited(time.Second)},
			{key: "c", at: 4 * time.Second, n: 1, want: nil},
		},
	}}
	for _, c := range cases {
		t.Run(c.behavior, func(t *testing.T) {
			t.Parallel()

			buckets := newBuckets(t, 2)
			for i, next := range c.calls {
				got := buckets.Take(next.key, start().Add(next.at), next.n)
				if !errors.Is(got, next.want) {
					t.Fatalf(
						"call %d: Take(%q, start+%v, %d) = %v, want %v",
						i+1, next.key, next.at, next.n, got, next.want,
					)
				}
			}
		})
	}

	t.Run("rejects n outside 1 to the capacity before looking for room", func(t *testing.T) {
		t.Parallel()

		buckets := newBuckets(t, 2)
		for _, key := range []string{"a", "b"} {
			err := buckets.Take(key, start(), 3)
			if err != nil {
				t.Fatalf("Take(%q, 3) = %v, want nil", key, err)
			}
		}

		for _, n := range []int{0, 4} {
			err := buckets.Take("c", start(), n)
			_, refused := errors.AsType[ratelimiting.LimitedError](err)
			if err == nil || refused {
				t.Errorf("Take(\"c\", %d) = %v, want an error other than LimitedError", n, err)
			}
		}
	})

	t.Run("admits concurrent callers no more new keys than its size", func(t *testing.T) {
		t.Parallel()

		buckets := newBuckets(t, 10)
		var (
			wg       sync.WaitGroup
			admitted atomic.Int64
		)
		for i := range 100 {
			wg.Go(func() {
				err := buckets.Take(strconv.Itoa(i), start(), 1)
				if err == nil {
					admitted.Add(1)
				}
			})
		}
		wg.Wait()

		if admitted.Load() != 10 {
			t.Errorf("%d of 100 new keys admitted, want 10", admitted.Load())
		}
	})
}

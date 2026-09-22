// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ratelimiting_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/andreychh/coopera/internal/ratelimiting"
)

func TestNewBucket(t *testing.T) {
	t.Parallel()

	t.Run("starts full", func(t *testing.T) {
		t.Parallel()

		bucket := ratelimiting.NewBucket(mustPolicy(t, 3))

		err := bucket.Take(start(), 3)
		if err != nil {
			t.Errorf("Take(3) on a new bucket of capacity 3 = %v, want nil", err)
		}
	})

	t.Run("rejects every Take when made from the zero Policy", func(t *testing.T) {
		t.Parallel()

		bucket := ratelimiting.NewBucket(ratelimiting.Policy{})

		err := bucket.Take(start(), 1)
		_, refused := errors.AsType[ratelimiting.LimitedError](err)
		if err == nil || refused {
			t.Errorf("Take(1) = %v, want an error other than LimitedError", err)
		}
	})
}

func TestBucket_Take(t *testing.T) {
	t.Parallel()

	// call is a Take made at start+at for n tokens, with the error it must
	// return.
	type call struct {
		at   time.Duration
		n    int
		want error
	}
	const ms = time.Millisecond
	drain := call{at: 0, n: 3, want: nil}

	// Every case plays its calls in order on a new bucket of capacity 3 that
	// earns one token per second.
	cases := []struct {
		behavior string
		calls    []call
	}{{
		behavior: "takes nothing when it refuses",
		calls: []call{
			{at: 0, n: 2, want: nil},
			{at: 0, n: 2, want: limited(time.Second)},
			{at: 0, n: 1, want: nil},
		},
	}, {
		behavior: "earns one token per interval",
		calls: []call{
			drain,
			{at: 999 * ms, n: 1, want: limited(ms)},
			{at: time.Second, n: 1, want: nil},
			{at: time.Second, n: 1, want: limited(time.Second)},
		},
	}, {
		behavior: "keeps the part of an interval that has not earned a token yet",
		calls: []call{
			drain,
			{at: 1400 * ms, n: 1, want: nil},
			{at: 2 * time.Second, n: 1, want: nil},
		},
	}, {
		behavior: "counts no time twice after a refusal",
		calls: []call{
			drain,
			{at: 1500 * ms, n: 2, want: limited(500 * ms)},
			{at: 1600 * ms, n: 2, want: limited(400 * ms)},
		},
	}, {
		behavior: "reports how long until n tokens are earned",
		calls: []call{
			drain,
			{at: 0, n: 3, want: limited(3 * time.Second)},
			{at: 400 * ms, n: 2, want: limited(1600 * ms)},
		},
	}, {
		behavior: "saves up no more than its capacity",
		calls: []call{
			drain,
			{at: time.Hour + 500*ms, n: 3, want: nil},
			{at: time.Hour + 500*ms, n: 1, want: limited(time.Second)},
		},
	}, {
		behavior: "treats a moment before the last one as the last one",
		calls: []call{
			{at: 10 * time.Second, n: 3, want: nil},
			{at: 0, n: 1, want: limited(time.Second)},
		},
	}}
	for _, c := range cases {
		t.Run(c.behavior, func(t *testing.T) {
			t.Parallel()

			bucket := ratelimiting.NewBucket(mustPolicy(t, 3))
			for i, next := range c.calls {
				got := bucket.Take(start().Add(next.at), next.n)
				if !errors.Is(got, next.want) {
					t.Fatalf(
						"call %d: Take(start+%v, %d) = %v, want %v",
						i+1, next.at, next.n, got, next.want,
					)
				}
			}
		})
	}

	t.Run("rejects n outside 1 to the capacity without taking anything", func(t *testing.T) {
		t.Parallel()

		bucket := ratelimiting.NewBucket(mustPolicy(t, 3))
		for _, n := range []int{0, -1, 4} {
			err := bucket.Take(start(), n)
			_, refused := errors.AsType[ratelimiting.LimitedError](err)
			if err == nil || refused {
				t.Errorf("Take(%d) = %v, want an error other than LimitedError", n, err)
			}
		}

		err := bucket.Take(start(), 3)
		if err != nil {
			t.Errorf("Take(3) after the rejected calls = %v, want nil", err)
		}
	})

	t.Run("gives concurrent callers no more than its capacity", func(t *testing.T) {
		t.Parallel()

		bucket := ratelimiting.NewBucket(mustPolicy(t, 100))
		var (
			wg    sync.WaitGroup
			taken atomic.Int64
		)
		for range 1000 {
			wg.Go(func() {
				err := bucket.Take(start(), 1)
				if err == nil {
					taken.Add(1)
				}
			})
		}
		wg.Wait()

		if taken.Load() != 100 {
			t.Errorf("%d of 1000 concurrent Take(1) succeeded, want 100", taken.Load())
		}
	})
}

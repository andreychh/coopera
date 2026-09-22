// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package ratelimiting limits how often an action may be taken.
package ratelimiting

import (
	"fmt"
	"sync"
	"time"
)

// LimitedError reports that an action was refused for now: a [Bucket] held
// too few tokens for it, or [Buckets] had no room for a new key. Wait is how
// soon after the refused moment the same action can succeed at the earliest;
// it may take longer if other actions get there first.
type LimitedError struct {
	Wait time.Duration
}

func (e LimitedError) Error() string {
	return fmt.Sprintf("rate limited, wait %v", e.Wait)
}

// Bucket is a token bucket that follows a [Policy].
//
// A Bucket is safe for concurrent use. The zero Bucket is not usable; create
// one with [NewBucket].
type Bucket struct {
	policy Policy

	mu sync.Mutex
	// level is the number of tokens the bucket held at countedAt.
	level int
	// countedAt advances by whole intervals, so time short of an interval
	// counts toward the next token. When the bucket fills up, countedAt is set
	// to the time of the call instead and the excess is dropped. It is zero
	// until a Take first counts tokens: a full bucket needs no reference
	// moment.
	countedAt time.Time
}

// NewBucket constructs a full Bucket that follows policy. A bucket made from
// the zero Policy returns an error from every Take.
func NewBucket(policy Policy) *Bucket {
	return &Bucket{
		policy:    policy,
		mu:        sync.Mutex{},
		level:     policy.capacity,
		countedAt: time.Time{},
	}
}

// Take removes n tokens from the bucket at the moment now. If the bucket
// holds fewer than n tokens, Take removes none and returns a [LimitedError].
// If n is less than 1 or greater than the policy's capacity, Take removes none
// and returns another error.
//
// Concurrent callers may pass moments out of order: a moment earlier than one
// already passed to Take adds no tokens.
func (b *Bucket) Take(now time.Time, n int) error {
	if !b.policy.admits(n) {
		return fmt.Errorf("the policy admits no action that costs %d tokens", n)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	capacity, interval := b.policy.capacity, b.policy.interval
	if now.Before(b.countedAt) {
		now = b.countedAt
	}
	elapsed := now.Sub(b.countedAt)
	if elapsed >= time.Duration(capacity-b.level)*interval {
		b.level = capacity
		b.countedAt = now
	} else {
		earned := int(elapsed / interval)
		b.level += earned
		b.countedAt = b.countedAt.Add(time.Duration(earned) * interval)
	}
	if n <= b.level {
		b.level -= n
		return nil
	}
	availableAt := b.countedAt.Add(time.Duration(n-b.level) * interval)
	return LimitedError{Wait: availableAt.Sub(now)}
}

// refillsAt returns the moment the bucket holds its whole capacity again if
// nothing more is taken from it. It is zero for a bucket never taken from.
func (b *Bucket) refillsAt() time.Time {
	b.mu.Lock()
	defer b.mu.Unlock()

	missing := b.policy.capacity - b.level
	return b.countedAt.Add(time.Duration(missing) * b.policy.interval)
}

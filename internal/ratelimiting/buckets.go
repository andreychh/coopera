// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ratelimiting

import (
	"fmt"
	"sync"
	"time"
)

// Buckets keeps a [Bucket] for each key, all following one [Policy], for at
// most size keys at a time. A key's bucket is dropped only once it has
// refilled, so a key that comes back finds the bucket it would have had
// anyway.
//
// Buckets is safe for concurrent use. The zero Buckets is not usable; create
// one with [NewBuckets].
type Buckets[K comparable] struct {
	policy Policy
	size   int

	mu      sync.Mutex
	buckets map[K]*Bucket
	// refillAt is no later than the moment any held bucket refills, so before
	// it a prune would drop nothing and a new key is refused without one. A
	// prune sets it to the earliest refill among the buckets it keeps; taking
	// from a bucket only delays its refill, and a new bucket lowers refillAt
	// to its own. The zero value calls for a prune.
	refillAt time.Time
}

// NewBuckets constructs an empty Buckets that holds at most size buckets
// following policy. It returns an error if size is less than 1.
func NewBuckets[K comparable](policy Policy, size int) (*Buckets[K], error) {
	if size < 1 {
		return nil, fmt.Errorf("size must be at least 1, got %d", size)
	}
	return &Buckets[K]{
		policy:   policy,
		size:     size,
		mu:       sync.Mutex{},
		buckets:  make(map[K]*Bucket),
		refillAt: time.Time{},
	}, nil
}

// Take removes n tokens from the bucket of key at the moment now, as
// [Bucket.Take] does; a key without a bucket gets a new, full one. While size
// buckets are held and none has refilled, Take refuses a new key with a
// [LimitedError]; a key that has a bucket is never refused for lack of room.
// If n is less than 1 or greater than the policy's capacity, Take returns
// another error and changes nothing.
func (b *Buckets[K]) Take(key K, now time.Time, n int) error {
	if !b.policy.admits(n) {
		return fmt.Errorf("the policy admits no action that costs %d tokens", n)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	bucket, found := b.buckets[key]
	if found {
		return bucket.Take(now, n)
	}
	if len(b.buckets) >= b.size && !now.Before(b.refillAt) {
		b.prune(now)
	}
	if len(b.buckets) >= b.size {
		return LimitedError{Wait: b.refillAt.Sub(now)}
	}
	bucket = NewBucket(b.policy)
	b.buckets[key] = bucket
	err := bucket.Take(now, n)
	at := bucket.refillsAt()
	if at.Before(b.refillAt) {
		b.refillAt = at
	}
	return err
}

// prune drops every bucket that is full at now and sets refillAt to the
// moment the first of the rest refills. Dropping a full bucket loses nothing:
// a new one starts full.
func (b *Buckets[K]) prune(now time.Time) {
	b.refillAt = time.Time{}
	for key, bucket := range b.buckets {
		at := bucket.refillsAt()
		if !at.After(now) {
			delete(b.buckets, key)
			continue
		}
		if b.refillAt.IsZero() || at.Before(b.refillAt) {
			b.refillAt = at
		}
	}
}

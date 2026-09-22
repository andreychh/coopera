// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ratelimiting

import (
	"fmt"
	"time"
)

// Policy is the pace a [Bucket] enforces: its capacity is the most tokens a
// bucket holds, and its interval is how long a bucket takes to earn one more.
//
// The zero Policy is not usable; create one with [NewPolicy].
type Policy struct {
	capacity int
	interval time.Duration
}

// NewPolicy constructs a Policy with the given capacity and interval. It
// returns an error if capacity is less than 1 or interval is not positive.
func NewPolicy(capacity int, interval time.Duration) (Policy, error) {
	if capacity < 1 {
		return Policy{}, fmt.Errorf("capacity must be at least 1, got %d", capacity)
	}
	if interval <= 0 {
		return Policy{}, fmt.Errorf("interval must be positive, got %v", interval)
	}
	return Policy{capacity: capacity, interval: interval}, nil
}

// admits reports whether a bucket following p can ever let through an action
// that costs n tokens.
func (p Policy) admits(n int) bool {
	return n >= 1 && n <= p.capacity
}

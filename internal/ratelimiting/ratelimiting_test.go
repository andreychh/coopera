// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ratelimiting_test

import (
	"testing"
	"time"

	"github.com/andreychh/coopera/internal/ratelimiting"
)

// start is the moment every test measures time from.
func start() time.Time {
	return time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
}

// mustPolicy returns a Policy with the given capacity and a one-second
// interval, and stops the test if NewPolicy rejects it.
func mustPolicy(t *testing.T, capacity int) ratelimiting.Policy {
	t.Helper()

	policy, err := ratelimiting.NewPolicy(capacity, time.Second)
	if err != nil {
		t.Fatalf("NewPolicy(%d, 1s) = %v, want nil", capacity, err)
	}
	return policy
}

// limited returns the error Take reports when it refuses an action that can
// succeed after wait at the earliest.
func limited(wait time.Duration) error {
	return ratelimiting.LimitedError{Wait: wait}
}

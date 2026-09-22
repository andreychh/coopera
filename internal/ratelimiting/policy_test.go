// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ratelimiting_test

import (
	"testing"
	"time"

	"github.com/andreychh/coopera/internal/ratelimiting"
)

func TestNewPolicy(t *testing.T) {
	t.Parallel()

	cases := []struct {
		behavior string
		capacity int
		interval time.Duration
		valid    bool
	}{{
		behavior: "accepts a positive capacity and interval",
		capacity: 1,
		interval: time.Nanosecond,
		valid:    true,
	}, {
		behavior: "rejects a capacity of zero",
		capacity: 0,
		interval: time.Second,
		valid:    false,
	}, {
		behavior: "rejects a negative capacity",
		capacity: -1,
		interval: time.Second,
		valid:    false,
	}, {
		behavior: "rejects an interval of zero",
		capacity: 1,
		interval: 0,
		valid:    false,
	}, {
		behavior: "rejects a negative interval",
		capacity: 1,
		interval: -time.Second,
		valid:    false,
	}}
	for _, c := range cases {
		t.Run(c.behavior, func(t *testing.T) {
			t.Parallel()

			_, err := ratelimiting.NewPolicy(c.capacity, c.interval)

			if (err == nil) != c.valid {
				t.Errorf(
					"NewPolicy(%d, %v) = %v, want valid = %t",
					c.capacity, c.interval, err, c.valid,
				)
			}
		})
	}
}

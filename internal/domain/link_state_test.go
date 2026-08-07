// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain_test

import (
	"testing"
	"time"

	"github.com/andreychh/coopera/internal/domain"
)

func TestNewLinkState(t *testing.T) {
	t.Parallel()

	at := func(hour int) domain.DateTime {
		return domain.DateTime(time.Date(2026, 8, 7, hour, 0, 0, 0, time.UTC))
	}
	past, now, future := at(11), at(12), at(13)

	cases := []struct {
		name      string
		expiresAt *domain.DateTime
		revokedAt *domain.DateTime
		want      domain.LinkState
	}{{
		name:      "no deadline and never revoked is active forever",
		expiresAt: nil,
		revokedAt: nil,
		want:      domain.Active{Validity: domain.Forever{}},
	}, {
		name:      "deadline ahead is active until it",
		expiresAt: &future,
		revokedAt: nil,
		want:      domain.Active{Validity: domain.Until{Time: future}},
	}, {
		name:      "deadline behind is expired",
		expiresAt: &past,
		revokedAt: nil,
		want:      domain.Expired{At: past},
	}, {
		// A link dies the moment its deadline arrives, not a moment after.
		name:      "deadline exactly now is expired",
		expiresAt: &now,
		revokedAt: nil,
		want:      domain.Expired{At: now},
	}, {
		name:      "revoked with no deadline is revoked",
		expiresAt: nil,
		revokedAt: &past,
		want:      domain.Revoked{At: past},
	}, {
		name:      "revoked before its deadline is revoked",
		expiresAt: &future,
		revokedAt: &past,
		want:      domain.Revoked{At: past},
	}, {
		// Revoking is something the owner did, expiring merely happened.
		name:      "revoked and expired reads as revoked",
		expiresAt: &past,
		revokedAt: &past,
		want:      domain.Revoked{At: past},
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := domain.NewLinkState(c.expiresAt, c.revokedAt, now)
			if got != c.want {
				t.Errorf("got %#v, want %#v", got, c.want)
			}
		})
	}
}

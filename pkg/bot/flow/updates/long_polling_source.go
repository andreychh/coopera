// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package updates

import (
	"context"

	"github.com/andreychh/coopera/pkg/bot/api"
	"github.com/andreychh/coopera/pkg/ptr"
)

type endpoint = Endpoint[api.GetUpdatesRequest, api.GetUpdatesResponse]

type LongPollingSource struct {
	endpoint endpoint
	buffer   int
}

func NewLongPollingSource(endpoint endpoint) LongPollingSource {
	return LongPollingSource{
		endpoint: endpoint,
		buffer:   10,
	}
}

func (s LongPollingSource) Updates(ctx context.Context) <-chan api.Update {
	channel := make(chan api.Update, s.buffer)
	go func() {
		defer close(channel)
		var offset int64
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			updates, err := s.endpoint.Call(
				ctx,
				api.GetUpdatesRequest{
					Offset:         ptr.Ptr(offset),
					Limit:          nil,
					Timeout:        ptr.Ptr(30),
					AllowedUpdates: nil,
				},
			)
			if err != nil {
				return
			}
			for _, update := range updates {
				offset = update.UpdateID + 1
				select {
				case <-ctx.Done():
					return
				case channel <- update:
				}
			}
		}
	}()
	return channel
}

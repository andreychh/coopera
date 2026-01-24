// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package updates

import (
	"context"

	"github.com/andreychh/coopera/pkg/bot/api"
	"github.com/andreychh/coopera/pkg/bot/endpoints"
	"github.com/andreychh/coopera/pkg/utils"
)

type endpoint = endpoints.Endpoint[api.GetUpdatesRequest, api.GetUpdatesResponse]

type longPollingSource struct {
	endpoint endpoint
	buffer   int
}

func (s longPollingSource) Updates(ctx context.Context) <-chan api.Update {
	channel := make(chan api.Update, s.buffer)
	go func() {
		defer close(channel)
		var offset int32
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			updates, err := s.endpoint.Call(
				ctx,
				api.GetUpdatesRequest{
					Offset:  utils.Ptr(offset),
					Timeout: utils.Ptr[int32](30),
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

func LongPollingSource(endpoint endpoint) UpdateSource {
	return longPollingSource{
		endpoint: endpoint,
	}
}

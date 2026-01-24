// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package endpoints

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera/pkg/bot/api"
	"github.com/andreychh/coopera/pkg/bot/transport"
)

type Endpoint[Request, Response any] interface {
	Call(ctx context.Context, req Request) (Response, error)
}

type endpoint[Request, Response any] struct {
	client transport.TelegramClient
	method api.Method
}

func (e endpoint[Request, Response]) Call(ctx context.Context, request Request) (Response, error) {
	var response Response
	err := e.client.SendRequest(ctx, e.method, request, &response)
	if err != nil {
		return response, fmt.Errorf("sending %s request: %w", e.method, err)
	}
	return response, nil
}

func GetMe(client transport.TelegramClient) Endpoint[api.GetMeRequest, api.GetMeResponse] {
	return endpoint[api.GetMeRequest, api.GetMeResponse]{
		client: client,
		method: api.MethodGetMe,
	}
}

func GetUpdates(client transport.TelegramClient) Endpoint[api.GetUpdatesRequest, api.GetUpdatesResponse] {
	return endpoint[api.GetUpdatesRequest, api.GetUpdatesResponse]{
		client: client,
		method: api.MethodGetUpdates,
	}
}

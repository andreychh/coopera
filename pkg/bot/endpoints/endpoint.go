// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package endpoints

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera/pkg/bot/api"
)

type Endpoint[Req, Resp any] struct {
	client TelegramClient
	method api.Method
}

//nolint:ireturn // Resp is a generic type parameter, not a functional interface return
func (e Endpoint[Req, Resp]) Call(ctx context.Context, req Req) (Resp, error) {
	var resp Resp
	err := e.client.SendRequest(ctx, e.method, req, &resp)
	if err != nil {
		return resp, fmt.Errorf("sending %s req: %w", e.method, err)
	}
	return resp, nil
}

func GetMe(client TelegramClient) Endpoint[api.GetMeRequest, api.GetMeResponse] {
	return Endpoint[api.GetMeRequest, api.GetMeResponse]{
		client: client,
		method: api.MethodGetMe,
	}
}

func GetUpdates(client TelegramClient) Endpoint[api.GetUpdatesRequest, api.GetUpdatesResponse] {
	return Endpoint[api.GetUpdatesRequest, api.GetUpdatesResponse]{
		client: client,
		method: api.MethodGetUpdates,
	}
}

func SendMessage(client TelegramClient) Endpoint[api.SendMessageRequest, api.SendMessageResponse] {
	return Endpoint[api.SendMessageRequest, api.SendMessageResponse]{
		client: client,
		method: api.MethodSendMessage,
	}
}

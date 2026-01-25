// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"encoding/json"
)

type Envelope struct {
	ErrorCode   *int32              `json:"error_code,omitempty"`
	Description *string             `json:"description,omitempty"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
	Result      json.RawMessage     `json:"result,omitempty"`
	Ok          bool                `json:"ok"`
}

type ResponseParameters struct {
	MigrateToChatId *int64 `json:"migrate_to_chat_id,omitempty"`
	RetryAfter      *int32 `json:"retry_after,omitempty"`
}

type GetMeRequest struct{}

type GetMeResponse User

type GetUpdatesRequest struct {
	Offset         *int32   `json:"offset,omitempty"`
	Limit          *int32   `json:"limit,omitempty"`
	Timeout        *int32   `json:"timeout,omitempty"`
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}

type GetUpdatesResponse []Update

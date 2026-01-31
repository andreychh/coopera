// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import "encoding/json"

type Envelope struct {
	ErrorCode   *int                `json:"error_code,omitempty"`
	Description *string             `json:"description,omitempty"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
	Result      json.RawMessage     `json:"result,omitempty"`
	Ok          bool                `json:"ok"`
}

type ResponseParameters struct {
	MigrateToChatId *int64 `json:"migrate_to_chat_id,omitempty"`
	RetryAfter      *int   `json:"retry_after,omitempty"`
}

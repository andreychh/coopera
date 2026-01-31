// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

type GetMeRequest struct{}

type GetMeResponse User

type GetUpdatesRequest struct {
	Offset         *int64   `json:"offset,omitempty"`
	Limit          *int     `json:"limit,omitempty"`
	Timeout        *int     `json:"timeout,omitempty"`
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}

type GetUpdatesResponse []Update

type SendMessageRequest struct {
	ChatID    ChatID    `json:"chat_id"`
	Text      string    `json:"text"`
	ParseMode ParseMode `json:"parse_mode,omitempty"`
	// todo: add missing fields.
}

type SendMessageResponse Message

// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

type Update struct {
	Message  *Message `json:"message"`
	UpdateID int32    `json:"update_id"`
}

type Message struct {
	MessageID int32 `json:"message_id"`
	Date      int32 `json:"date"`
}

type User struct {
	Username  *string `json:"username,omitempty"`
	FirstName string  `json:"first_name"`
	ID        int64   `json:"id"`
	IsBot     bool    `json:"is_bot"`
}

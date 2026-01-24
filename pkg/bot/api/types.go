// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

type Update struct {
	UpdateID int32    `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	MessageID int32 `json:"message_id"`
	Date      int32 `json:"date"`
}

type User struct {
	ID        int64   `json:"id"`
	IsBot     bool    `json:"is_bot"`
	FirstName string  `json:"first_name"`
	Username  *string `json:"username,omitempty"`
}

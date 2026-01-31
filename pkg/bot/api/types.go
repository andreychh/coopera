// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

type Update struct {
	Message  *Message `json:"message,omitempty"`
	UpdateID int64    `json:"update_id"`
}

type Message struct {
	From      *User `json:"from,omitempty"`
	MessageID int64 `json:"message_id"`
	Date      int64 `json:"date"`
}

type User struct {
	Username  *string `json:"username,omitempty"`
	FirstName string  `json:"first_name"`
	ID        int64   `json:"id"`
	IsBot     bool    `json:"is_bot"`
}

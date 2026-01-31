// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"encoding/json"
	"errors"
)

// ErrInvalidChatID is returned when a ChatID struct has no fields set.
var ErrInvalidChatID = errors.New("ChatID must have either ChatID or ChannelUsername set")

type ChatID struct {
	ChatID          *int
	ChannelUsername *string
}

func (id ChatID) MarshalJSON() ([]byte, error) {
	if id.ChatID != nil {
		return json.Marshal(id.ChatID)
	}
	if id.ChannelUsername != nil {
		return json.Marshal(id.ChannelUsername)
	}
	return nil, ErrInvalidChatID
}

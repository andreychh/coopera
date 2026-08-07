// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import "github.com/google/uuid"

// ID names one thing in the system. Which thing is told by where the ID
// sits, never by the value itself.
type ID uuid.UUID

func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}

func (i ID) String() string {
	return uuid.UUID(i).String()
}

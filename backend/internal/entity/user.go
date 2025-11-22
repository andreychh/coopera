package entity

import "time"

type UserEntity struct {
	ID         *int32
	TelegramID *int64
	Username   *string
	CreatedAt  *time.Time
	Teams      []TeamWithRole // ← сюда добавляем информацию о командах
}

type TeamWithRole struct {
	TeamID   int32
	TeamName string
	Role     Role
}

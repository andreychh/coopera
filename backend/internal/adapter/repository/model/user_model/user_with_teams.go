package user_model

import "time"

type UserWithTeams struct {
	ID         int32
	TelegramID int64
	Username   string
	CreatedAt  time.Time
	Teams      []TeamWithRole
}

type TeamWithRole struct {
	ID   int32
	Name string
	Role string
}

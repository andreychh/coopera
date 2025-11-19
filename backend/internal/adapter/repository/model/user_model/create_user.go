package user_model

import "time"

type User struct {
	ID         int32     `db:"id"`
	TelegramID int64     `db:"telegram_id"`
	CreatedAt  time.Time `db:"created_at"`
}

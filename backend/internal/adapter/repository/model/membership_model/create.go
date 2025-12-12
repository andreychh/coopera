package membership_model

import "time"

type Membership struct {
	ID        int32     `db:"id"`
	TeamID    int32     `db:"team_id"`
	UserID    int32     `db:"user_id"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}

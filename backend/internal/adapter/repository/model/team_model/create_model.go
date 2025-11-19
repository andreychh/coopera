package team_model

import "time"

type Team struct {
	ID        int32     `db:"id"`
	Name      string    `db:"name"`
	CreatedBy int32     `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
}

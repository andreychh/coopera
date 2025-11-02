package entity

import "time"

type TeamEntity struct {
	ID        *int32
	Name      string
	CreatedAt *time.Time
	CreatedBy int32
}

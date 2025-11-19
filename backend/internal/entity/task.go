package entity

import "time"

type Task struct {
	ID          int32
	TeamID      int32
	Title       string
	Description *string
	Points      int32
	Status      *Status
	AssignedTo  *int32
	CreatedBy   int32
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
